# Lavaboom API

To install:

```
go get github.com/lavab/api
api
curl --data "username=abc&password=def" localhost:5000/signup
curl --data "username=abc&password=def" localhost:5000/login
curl --header "Auth: <token>" localhost:5000/me
```

## Configuration variables

You can use either commandline flags:
```
{ api } master » ./api -help
Usage of api:
  -bind=":5000": Network address used to bind
  -classic_registration=false: Classic registration enabled?
  -force_colors=false: Force colored prompt?
  -log="text": Log formatter type. Either "json" or "text"
  -rethinkdb_db="dev": Database name on the RethinkDB server
  -rethinkdb_key="": Authentication key of the RethinkDB database
  -rethinkdb_url="localhost:28015": Address of the RethinkDB database
  -session_duration=72: Session duration expressed in hours
  -version="v0": Shown API version
```

Or environment variables:
```
{ api } master » BIND=127.0.0.1:5000 CLASSIC_REGISTRATION=false \
FORCE_COLORS=false LOG=text RETHINKDB_DB=dev RETHINKDB_KEY="" \
RETHINKDB_URL=localhost:28015 SESSION_DURATION=72 VERSION=v0 ./api
```

Or configuration files:
```
{ api } master » cat api.conf
# lines beggining with a "#" character are treated as comments
bind :5000
classic_registration false
force_colors false
log text

rethinkdb_db dev
# configuration values can be empty
rethinkdb_key
# Keys and values can be also seperated by the "=" character
rethinkdb_url=localhost:28015

session_duration=72
version=v0
{ api } master » ./api -config api.conf
```

## Build status:

 - `master` - [![Circle CI](https://circleci.com/gh/lavab/api/tree/master.svg?style=svg&circle-token=4a52d619a03d0249906195d6447ceb60a475c0c5)](https://circleci.com/gh/lavab/api/tree/master)
 - `develop` - [![Circle CI](https://circleci.com/gh/lavab/api/tree/develop.svg?style=svg&circle-token=4a52d619a03d0249906195d6447ceb60a475c0c5)](https://circleci.com/gh/lavab/api/tree/develop)
