# Dockerfile used to create "husyci/enry" image
# https://hub.docker.com/r/huskyci/enry/

FROM golang:alpine as builder

RUN apk update && apk upgrade \
	&& apk add --no-cache alpine-sdk git \
	&& apk --update add jq

RUN go get gopkg.in/src-d/enry.v1/...

# From the base image
FROM alpine:3.8
COPY --from=builder /go/bin/enry /usr/bin/
COPY --from=builder /usr/bin/jq /usr/bin/
