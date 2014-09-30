FROM google/golang

RUN go get github.com/lavab/api

CMD []
ENTRYPOINT ["/gopath/bin/api"]
