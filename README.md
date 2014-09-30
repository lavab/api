# Lavaboom API

To install:

```
go get github.com/lavab/api
sudo api
curl --data "username=abc&password=def" localhost/signup
curl --data "username=abc&password=def" localhost/login
curl --data "token=???" --get localhost/me                  # or -X GET
```
