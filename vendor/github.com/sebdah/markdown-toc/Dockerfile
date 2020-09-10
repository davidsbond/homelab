FROM golang:1.9-stretch

RUN mkdir -p /go/src/github.com/sebdah/markdown-toc
WORKDIR /go/src/github.com/sebdah/markdown-toc
ADD . /go/src/github.com/sebdah/markdown-toc

RUN make install

ENTRYPOINT ["markdown-toc"]
