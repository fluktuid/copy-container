# Copy-Container

This project contains a simple go service for copying from one to another folder.

The copying is especially useful if you need a tiny distroless container copying data from one pvc to another, e.g. in a kubernetes cluster.

## Configuration

The App configuration is done via env variables.

| env_var            | description                                      | default |
| ------------------ | ------------------------------------------------ | ------- |
| FROM               | folder to copy from                              | /from   |
| TO                 | folder to copy to                                | /to     |
| SET_PERMISSION     | if set, overrides permission of all copied files | -       |
| SYNC               | if true, sync folder after copy                  | true    |
| SKIP_ON_EXIST      | if true, skips existing subfolder while copying  | true    |
| PRESERVE_TIMES     | preserver stat times while copying               | true    |
| BUFFER             | copy buffer (default: 32k)                       | 32768   |
| LOG_FORMATTER_JSON | if true, log output is formatted to json         | false   |
| LOG_LEVEL          | used log level                                   | warn    |

## Build

``` bash
$ docker build .
```

### Dependencies

- [otiai10/copy](github.com/otiai10/copy) (MIT)
- [kelseyhightower/envconfig](github.com/kelseyhightower/envconfig) (MIT)
- [sirupsen/logrus](github.com/sirupsen/logrus) (MIT)
