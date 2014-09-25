# Lavaboom API

This is a draft of Lavaboom's back-end API.

## HTTP calls

```
Public calls

GET     /                           information about the API (version, etc)

POST    /login                      sign in into existing account
POST    /signup                     registers a new account
GET     /users/:name/exists         checks whether a user exist
GET     /keys/:id                   id can be fingerprint, email, or key id


Session calls

GET     /me                         returns the name of the current user
POST    /logout                     invalidates current token

GET     /users/:name/messages       gets all the messages, newest first
PUT     /users/:name                updates user information
PUT     /users/:name/password
DELETE  /users/:name

GET     /threads
GET     /threads/:id

POST    /messages                   new message
GET     /messages/:id               gets individual message
DELETE  /messages/:id
PUT     /messages/:id

GET     /tags
GET     /tags/:name
GET     /tags/:name/emails

GET     /contacts
GET     /contacts/:id
POST    /contacts/:id
PUT     /contacts/:id
GET     /contacts/:id/threads

POST    /keys
```

## Data models

The default data type is String.

```
user:
  id
  name
  namePref
  password
  salt
  pgp:
    key
    finger
    algo
    expiry            Date
  app:
    domainPref
    timezone
    #etc
  billing:
    name
    address
    country
    birth             Date
    method
    #etc

message:
  id
  from
  thread                                    thread ID
  sent                Date
  tags                String[]
  #etc

thread:
  id
  members
  encrypted           Boolean
  messages[]:
    id
    preview
    timestamp         Date
  #etc

file:
  id
  name
  mime
  size                Int64               size in bytes (?)
  payload
  crypto
    key
    method

tag:
  id
  name
  immutable           Bool                think Inbox, etc.
  icon
  count
    all
    unread

contact:
  id
  name
  email
  avatar
  pgp
    key
    finger
    id
  #etc
```
