FROM google/golang:1.4

RUN go get github.com/tools/godep

RUN mkdir -p /gopath/src/github.com/lavab/api
ADD . /gopath/src/github.com/lavab/api
RUN cd /gopath/src/github.com/lavab/api && godep go install

ENTRYPOINT ["/gopath/bin/api"]
