FROM golang:alpine as builder
RUN apk update && apk upgrade

# Install required tools
RUN apk add --no-cache ca-certificates make bash git

ADD . /project
WORKDIR /project

# Compile binaries
RUN make

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /project/bin /bin
