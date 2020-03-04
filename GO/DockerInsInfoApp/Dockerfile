# build stag
FROM golang:alpine
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN go get github.com/go-redis/redis && go get -u github.com/go-sql-driver/mysql
ADD . /go/src
RUN cd /go/src && go build -o MyApp
WORKDIR /go/src
ENTRYPOINT ["./MyApp", "1>/dev/stdout", "2>/dev/stderr"]
