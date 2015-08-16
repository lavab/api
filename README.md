# Lavaboom API

[![Code Climate](https://codeclimate.com/github/lavab/api/badges/gpa.svg)](https://codeclimate.com/github/lavab/api)

<img src="https://mail.lavaboom.com/img/Lavaboom-logo.svg" align="right" width="200px" />

Lavaboom's main API written in Golang.

Contains the core Lavaboom's functionality. Currently a monolith, which later
will be split and transformed into a set of microservices. 

## Requirements

 - Redis
 - RethinkDB
 - NSQ

## Installation

### Inside a Docker container

*This image will be soon uploaded to Docker Hub.*

```bash
git clone https://github.com/lavab/api.git
cd api
docker build -t "lavab/api" .
docker run \
	-p 127.0.0.1:5000:5000 \
	--name api \
	lavab/api \
	-redis_address=172.8.0.1:6379 \
	-lookupd_address=172.8.0.1:4161 \
	-nsqd_address=172.8.0.1:4150 \
	-rethinkdb_address=172.8.0.1:28015 \
	-api_host=api.lavaboom.com \
	-email_domain=lavaboom.com
```

### Directly running the parts

```bash
go get github.com/lavab/api

api \
	-redis_address=172.8.0.1:6379 \
	-lookupd_address=172.8.0.1:4161 \
	-nsqd_address=172.8.0.1:4150 \
	-rethinkdb_address=172.8.0.1:28015 \
	-api_host=api.lavaboom.com \
	-email_domain=lavaboom.com
```

## Passing configuration

You can use either commandline flags:
```
{ api } master » ./api -help
Usage of api:
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
  -lookupd_address="127.0.0.1:4160": Address of the lookupd server
  -nsqd_address="127.0.0.1:4150": Address of the nsqd server
  -redis_address="127.0.0.1:6379": Address of the redis server
  -redis_db=0: Index of redis database to use
  -redis_password="": Password of the redis server
  -rethinkdb_address="127.0.0.1:28015": Address of the RethinkDB database
  -rethinkdb_db="dev": Database name on the RethinkDB server
  -rethinkdb_key="": Authentication key of the RethinkDB database
  -session_duration=72: Session duration expressed in hours
  -slack_channel="#notif-api-logs": channel to which Slack bot will send messages
  -slack_icon=":ghost:": emoji icon of the Slack bot
  -slack_level="warning": minimal level required to have messages sent to slack
  -slack_url="": URL of the Slack Incoming webhook
  -slack_username="API": username of the Slack bot
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

## License

This project is licensed under the MIT license. Check `license` for more
information.
