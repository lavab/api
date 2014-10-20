# Lavaboom API

To install:

```
go get github.com/lavab/api
api
curl --data "username=abc&password=def" localhost:5000/signup
curl --data "username=abc&password=def" localhost:5000/login
curl --header "Auth: <token>" localhost:5000/me
```
