FROM golang:1.13-alpine3.11 as builder
# All these steps will be cached
RUN apk add --update --no-cache ca-certificates git
ENV GO111MODULE=on
WORKDIR /app
RUN go env -w GOPRIVATE=g.ghn.vn/logistic/*
#RUN go env -w GOPROXY=https://goproxy.io,direct
RUN git config --global \
  url."https://tientp:qH8uisEJZZNrxxcdwy-c@g.ghn.vn/".insteadOf \
  "https://g.ghn.vn/"
RUN go get -u -v github.com/tmc/pqstream/cmd/...
# COPY the source code as the last step
COPY go.mod .
COPY go.sum .
COPY  . .
RUN ls
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix cgo -o /go/bin/app_build main.go
RUN echo $GOPATH
RUN ls /go/bin/
# Run docker entrypoint
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 7000
# CMD ["pqsd"]
