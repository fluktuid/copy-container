# based on https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
# STEP 1 build executable binary
############################
FROM golang:1.17-alpine AS builder
ARG GO_FILES=.
ARG GO_MAIN=main.go
LABEL maintainer="Lukas Paluch <fluktuid@users.noreply.github.com>"
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Create appuser.
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Fetch dependencies.
COPY ${GO_FILES}/go.mod .
COPY ${GO_FILES}/go.sum .
RUN go mod download
# Using go get.
RUN go get -v all

# get go files
COPY ${GO_FILES} .
# Build the binary.
ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=0
RUN GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-w -s" -o /app/main ${GO_MAIN}

# create tmp folder for use in scratch
RUN mkdir /my_tmp
RUN chown -R ${USER}:${USER} /my_tmp


############################
# STEP 2 build a small image
############################
FROM scratch
LABEL maintainer="Lukas Paluch <fluktuid@users.noreply.github.com>"
ARG PROJECT_NAME=copy-container
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# create tmp directory
COPY --from=builder /my_tmp /tmp
# Copy static executable.
COPY --from=builder /app/main /main

# Use an unprivileged user.
USER ${USER}:${USER}


# Run the binary.
ENTRYPOINT ["/main"]

