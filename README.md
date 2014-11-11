# Lavaboom API

To install:

```
go get github.com/lavab/api
api
curl --data "username=abc&password=def" localhost:5000/signup
curl --data "username=abc&password=def" localhost:5000/login
curl --header "Auth: <token>" localhost:5000/me
```

## Build status:

 - `master` - [![Build Status](https://magnum.travis-ci.com/lavab/api.svg?token=kJbppXeTxzqpCVvt4t5X&branch=master)](https://magnum.travis-ci.com/lavab/api)
 - `develop` - [![Build Status](https://magnum.travis-ci.com/lavab/api.svg?token=kJbppXeTxzqpCVvt4t5X&branch=develop)](https://magnum.travis-ci.com/lavab/api)