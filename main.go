package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	From          string `envconfig:"FROM" default:"/from"`
	To            string `envconfig:"TO" default:"/to"`
	SetPermission uint16 `envconfig:"SET_PERMISSION"`
	Sync          bool   `envconfig:"SYNC" default:"true"`
	SkipOnExist   bool   `envconfig:"SKIP_ON_EXIST" default:"true"`
	PreserveTimes bool   `envconfig:"PRESERVE_TIMES" default:"true"`
	Buffer        uint   `envconfig:"BUFFER" default:"32768"`
	LogJson       bool   `envconfig:"LOG_FORMATTER_JSON" default:"false"`
	LogLevel      string `envconfig:"LOG_LEVEL" default:"warn"`
}

func main() {
	var c Config
	err := envconfig.Process("", &c)
	setLogging(c.LogJson, c.LogLevel)
	if err != nil {
		log.Fatal(err.Error())
	}
	opt := copy.Options{
		Sync:           c.Sync,
		PreserveTimes:  c.PreserveTimes,
		CopyBufferSize: c.Buffer,
	}
	if c.SkipOnExist {
		opt.OnDirExists = func(_, dest string) copy.DirExistsAction {
			if dest != "from" {
				return copy.Untouchable
			}
			return copy.Merge
		}
	}
	if c.SetPermission > 0 {
		if c.SetPermission > 777 {
			c.SetPermission = 777
		}
		opt.AddPermission = os.FileMode(int(c.SetPermission))
	}

	err = copy.Copy(c.From, c.To, opt)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setLogging(jsonFormat bool, level string) {
	if jsonFormat {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if level != "" {
		ll, err := log.ParseLevel(level)
		if err != nil {
			log.Warn("can't parse log level, using default")
		} else {
			log.WithField("level", level).
				Info("set global log level")
		}
		log.SetLevel(ll)
	}
}
