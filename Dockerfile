FROM gliderlabs/alpine:latest

RUN apk update && apk add ca-certificates && apk add curl

ADD fuelsales-rpt /go/src/
COPY config/*.yaml /go/src/config/

WORKDIR /go/src

ENTRYPOINT /go/src/fuelsales-rpt