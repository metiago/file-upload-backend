FROM golang:1.11-alpine as base
RUN apk add --no-cache libstdc++ gcc g++ make git ca-certificates linux-headers
LABEL Tiago Souza
WORKDIR /go/src/github.com/metiago/zbx1
ADD . .
ENV GO111MODULE=on
RUN go get && go mod tidy

FROM alpine:latest
LABEL Tiago Souza
RUN apk add --no-cache jq ca-certificates linux-headers
COPY --from=base /go/bin/zbx1 /usr/local/bin/zbx1


EXPOSE 5000
ENTRYPOINT zbx1
