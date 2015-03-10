[![Stories in Ready](https://badge.waffle.io/lavab/api.png?label=ready&title=Ready)](https://waffle.io/lavab/api)
# Lavaboom API

To install:

```
go get github.com/lavab/api
./api
curl localhost:5000
```

## Configuration variables

You can use either commandline flags:
```
{ api } master » ./api -help
Usage of api
  -api_version="v0": Shown API version
  -bind=":5000": Network address used to bind
  -config="": config file to load
  -email_domain="lavaboom.io": Domain of the default email service
  -etcd_address="": etcd peer addresses split by commas
  -etcd_ca_file="": etcd path to server cert's ca
  -etcd_cert_file="": etcd path to client cert file
  -etcd_key_file="": etcd path to client key file
  -etcd_path="settings/": Path of the keys
  -force_colors=false: Force colored prompt?
  -log="text": Log formatter type. Either "json" or "text"
  -loggly_token="": Loggly token
  -nats_address="nats://127.0.0.1:4222": Address of the NATS server
  -redis_address="127.0.0.1:6379": Address of the redis server
  -redis_db=0: Index of redis database to use
  -redis_password="": Password of the redis server
  -rethinkdb_address="127.0.0.1:28015": Address of the RethinkDB database
  -rethinkdb_db="dev": Database name on the RethinkDB server
  -rethinkdb_key="": Authentication key of the RethinkDB database
  -session_duration=72: Session duration expressed in hours
  -yubicloud_id="": YubiCloud API id
  -yubicloud_key="": YubiCloud API key
```

Or environment variables:
```
{ api } master » BIND=127.0.0.1:6000 VERSION=v1 ./api
```

Or configuration files:
```
{ api } master » cat api.conf
# lines beggining with a "#" character are treated as comments
bind :5000
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

## Flags list

| Flag or config variable name | Environment variable name | Default value | Description |
|:---------------------------- | ------------------------- | --------| ----- |
| api_version    | `API_VERSION` | v0 | Shown API version |
| bind           | `BIND` | 0.0.0.0:5000 | Network address used to bind |
| config         | `CONFIG` | _empty_ | Config file to load |
| email_domain   | `EMAIL_DOMAIN` | lavaboom.io | Domain of the default email service |
| etcd_address   | `ETCD_ADDRESS` | _empty_ | etcd peer addresses split by commas |
| etcd_ca_file   | `ETCD_CA_FILE` | _empty_ | etcd path to server cert's ca |
| etcd_cert_file | `ETCD_CERT_FILE` | _empty_ | etcd path to client cert file |
| etcd_key_file  | `ETCD_KEY_FILE` | _empty_ |  etcd path to client key file |
| etcd_path      | `ETCD_PATH` | _empty_ | Path of the keys |
| force_colors   | `FORCE_COLORS` | false | Force colored prompt? |
| log            | `LOG` | text | Log formatter type. Either "json" or "text". |
| loggly_token   | `LOGGLY_TOKEN` | _empty_ |  Token used to connect to Loggly |
| nats_address   | `NATS_ADDRESS` | nats://127.0.0.1:4222 | Address of the NATS server |
| redis_address  | `REDIS_ADDRESS` | 127.0.0.1:6379 | Address of the redis server |
| redis_db       | `REDIS_DB` | 0 | Index of redis database to use |
| redis_password | `REDIS_PASSWORD` | _empty_ | Password of the redis server |
| rethinkdb_address | `RETHINKDB_ADDRESS` | 127.0.0.1:28015 | Address of the RethinkDB database |
| rethinkdb_db      | `RETHINKDB_DB` | dev |Database name on the RethinkDB server |
| rethinkdb_key     | `RETHINKDB_KEY` | _empty_ | Authentication key of the RethinkDB database |
| session_duration  | `SESSION_DURATION` | 72 | Session duration expressed in hours |
| yubicloud_id      | `YUBICLOUD_ID` | _empty_ | YubiCloud API ID. |
| yubicloud_key     | `YUBICLOUD_KEY` | _empty_ | YubiCloud API key. |

## Build status:

 - `master` - [![Circle CI](https://circleci.com/gh/lavab/api/tree/master.svg?style=svg&circle-token=4a52d619a03d0249906195d6447ceb60a475c0c5)](https://circleci.com/gh/lavab/api/tree/master) [![Coverage Status](https://coveralls.io/repos/lavab/api/badge.svg?branch=master)](https://coveralls.io/r/lavab/api?branch=master)
 - `develop` - [![Circle CI](https://circleci.com/gh/lavab/api/tree/develop.svg?style=svg&circle-token=4a52d619a03d0249906195d6447ceb60a475c0c5)](https://circleci.com/gh/lavab/api/tree/develop) [![Coverage Status](https://coveralls.io/repos/lavab/api/badge.svg?branch=develop)](https://coveralls.io/r/lavab/api?branch=develop)
