FROM golang:alpine as builder

# Install required tools
RUN apk add --update --no-cache ca-certificates make bash git upx

ADD . /project
WORKDIR /project

# Compile binaries
RUN make
RUN make pack

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /project/bin /bin
