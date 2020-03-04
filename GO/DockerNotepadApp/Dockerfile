# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN go get github.com/santosh07bec/MyDockerGoWebApp
RUN cd /go/src/github.com/santosh07bec/MyDockerGoWebApp && GOVENDOREXPERIMENT=1 go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/santosh07bec/MyDockerGoWebApp /app/
ENTRYPOINT ["./goapp", "1>/dev/stdout", "2>/dev/stderr"]

