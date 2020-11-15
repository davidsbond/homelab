FROM golang:alpine as builder

# Install required tools
RUN apk add --update --no-cache ca-certificates make bash git

ADD . /project
WORKDIR /project

# Compile binaries
RUN make

# Currently, upx is not available as an alpine package for arm64 devices:
# https://github.com/upx/upx/issues/419
# RUN make pack

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /project/bin /bin
