FROM google/golang

MAINTAINER "Andrei Simionescu <andrei@lavaboom.com>"

RUN mkdir -p /gopath/src/github.com/lavab/api
ADD . /gopath/src/github.com/lavab/api
#RUN go get github.com/lavab/api

RUN cd /gopath/src/github.com/lavab/api && godep go install

CMD []
ENTRYPOINT ["/gopath/bin/api"]
