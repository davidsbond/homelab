FROM golang:1.17.6 as builder

# Install required tools
RUN apt-get update -y
RUN apt-get install -y upx

ADD . /project
WORKDIR /project

# Compile binaries
RUN make
RUN make pack

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /project/bin /bin
