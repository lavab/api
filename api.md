# API

## What?

A documentation  of the first version (v0) of the Lavaboom API, focused on typical user interactions with the client.

### Clarification

* `client` is any Angular, iOS, Android, etc application that uses the Lavaboom API
* `user` is an actual person using a Lavaboom app
* `TBP` To Be Planned, i.e. not a v0 feature

## Most common actions

### Composing an email

```
POST /emails
Auth: {token}

{
	"to": ["example@domain.com","Andrei Simionescu", "andrei@lavaboom.com"],
	"cc": ["{some_user_id}"],
	"bcc": [],
	"reply_to": "{email_id}",
	"thread_id": "{thread_id}",
	"body": "{email_text}",
	"attachment_ids": ["{some_uuid}, {other_uuid}"]
}
```

1. `token` is the auth token, passed as an `Auth` HTTP header.
2. `to`, `cc`, `bcc` are arrays of strings that can be an email address, a full name, or a user ID. They are all optional, for instance if `thread_id` is specified and no new members are introduced to the discussion.
3. `email_id` if this is specified we know how this email fits in a current email discussion.
4. `thread_id` if this is a reply to an existing thread.
5. `email_text` is the content of the email, excluding attachments.
6. `attachment_ids` is a list of attachment IDs. Attachment uploading: POST to `/upload`, returns `attachment_id`s.
If they're not used, uploaded attachments expire in 24h. Also, while will in beta, the upload limit per file will be 25MB.
The attachments must be PGP-encrypted, so the user must be warned that changing the list of recipients after uploading an
attachment means he will have to re-upload that file.

Response:

1. `success` can be `delivered` or `error`
2. `email_id` is the ID of the new email, if the operation was successful
3. `message` is usually an error code, if necessary

Issues:

1. This can be trickier if we add **drafts** and **autosave** functionality, so we'll skip that for v0.
2. By adding websockets we can check for successful delivery and notify the user that the email has been delivered.
This means that the API will return `success: "queued"` or similar when POST-ing to `/emails`.
3. If this is an encrypted email (detectable on the server side automatically), the client must have already figured
out the public keys and performed encryption.

### Checking the inbox

```
GET /threads?label=inbox
Auth: {token}
```
1. `token` is the auth token
2. `offset` and `limit`, not specified, handle pagination and are by default 0 and 50.
3. `label` can be a label name or a label ID. Querying using a ID is faster.

Response:
```
{
	"n_threads": {number_of_threads},
	"threads": [
		{
			"thread_id": "{thread_id}",
			"participants": "{participants_snippet}",
			"subject": "{subject}",
			"snippet": "{snippet}",
			"labels": ["inbox","{some_label}"],
			"n_emails": {n_emails},
			"has_attachments": {has_attachments},
			"is_starred": {is_starred},
			"encription_level": "{encryption_level}",
			"avatar": "{avatar_id}",
			"is_postponed": {is_postponed}
		},
		[...]
	]
}
```

1. `number_of_threads` is an integer representing the number of threads found.
2. `thead_id` is provided because more data about the thread can be extracted by querying explicitly.
3. `participants_snippet` is a string that indicates to the user who the last interaction happened with.
Example: "Billy Joel", "To: John Smith, Billy Joel, and 256 others", etc.
4. `subject` is the subject of the latest email in the thread.
5. `snippet` is a snippet of the latest email in the thread.
6. `some_label` is a label name. It would be pointless to return UUIDs here.
7. `n_emails` is the number of emails in this thread. No need to report the number of unread emails separately.
8. `has_attachments` is a bool value indicating whether the last email has any attachments. (Or the whole
thread? This is not decided yet.)
9. `is_starred` is a bool value indicating whether this thread has been starred.
10. `encryption_level` can be a couple of string values (unencrypted, strong, weak, medium, and danger) and
depends on how many emails in the thread are encrypted and on the quality of the public keys used for encryption.
11. `avatar_id` will be an id for an avatar (of the person last interacted with in this thread). TBP.
12. `is_postponed` will be useful when we implement Mailbox-like productivity features.

Issues:

1. Add websockets to push new emails to the client without an explicit query.

### Searching

API v0 will only feature simple keyword search. The smart and [fuzzy](https://www.google.de/search?q=fuzzy+search) bits will happen server side.

```
GET /threads?q=my+search+query
Auth: {token}
```

You can combine label and keyword search.

```
GET /threads?label=inbox&q=where+my+email+at
Auth: {token}
```

The data can be sent via JSON too, although I think the URL way is more common for GET requests. The response is identical to the previous action.

## Other user actions

### Read emails for a specific label 

```
GET /threads?label={label}
Auth: {token}
```

Response: See "Checking the inbox"

### User settings

```
GET /me/settings
Auth: {token}
```

```
PUT /me/settings
Auth: {token}

{
	"{key}": "{value}"
}

```

Pretty self explanatory. [check out possible values [here](https://developer.zendesk.com/rest_api/docs/core/account_settings) or similar]

### Managing emails

Typical actions on emails that the user performs in their inbox:
* archive
* reply
* report spam
* mark unread
* (TBP) postpone/reschedule email
* (TBP) mute thread
* delete
* move

Issues:

1. A simple mechanism for **batching requests** needs to be introduced to handle
such tasks. Usually multiple emails are handled at once (archive, postpone, etc).

### Contacts

#### Adding a contact

```
POST /contacts
Auth: {token}

{
	"name": "John Smith",
	"email": {
		"address": "example@domain.com",
		"public_key": "{pgp_key}",
		"fingerprint": "{fingerprint}"
	}
	"fields": {
		"first_name": "John",
		"family_name": "Smith",
		"iphone": "0044123456789",
		"{key}": "{value}"
	}
}
```

1. `token` is the authentication token
2. `pgp_key` is the full PGP key for this email address. Any eventual ambiguity will be solved client-side.
3. `fingerprint` can be submitted together with the key, for added security. If only the fingerprint is submitted, the server will try to match it with a key on the Lavaboom key server, or if that fails, on a public key server.

Response:

```
{
	"success": "{success}",
	"message": "{message}"
}
```

1. `success` is a boolean. (Should it be a string?)
2. `message` explains what happened in case of an error. For example, an email address could've already been in the contacts' list.

#### Searching contacts

```
GET /contacts/?q=john+smith
Auth: {token}
```

Response: similar to what's submitted when adding a contact.

#### Deleting a contact

```
DELETE /contacts/{id}
Auth: {token}
```

Response: `success` and `message`

#### Adding a public key (or simply editing a contact)

```
PUT /contacts/{id}
Auth: {token}

{
	"{key}": "{value}"
}
```

#### Discussion about server public keys management

The user interaction with the Contacts screen is pretty straightforward, the only interesting aspect is added
by the fact that we need to manage public keys too. There are a number of ways to approach this, we must (as always)
try to balance the utmost security with convenience.

The options are:

1. completely manual key management; additionally we might let user sync with reputable key servers
2. completely managed key management; the Lavaboom key server is considered the ultimate authority on public key authenticity
3. flexible managed scheme; vote-based system – Lavaboom's key server can have multiple user-sourced public keys per email address,
and the users can vote on their authenticity, and elect which key to use (or to upload their own)

I personally prefer the last option, eventually augmented by using the existing public key servers to populate our key server with the keys corresponding to the email addresses used by our users (obviously, in a batched, unobvious way, to protect our users' privacy).

## Sessions

### Signup

```
POST /signup

{
	"username": "{user}",
	"password": "{pass}",
	"reg_token": "{reg}"
}
```

1. `user` is the desired username, {a-zA-Z0-9_.}, length between 2 and 32 characters. Like Gmail, the actual username doesn't contain periods, but the user has one preferred style (e.g. `john.smith`), and can use any variation of his username with periods, as long as they're not subsequent characters or at the beginning or the end of the username.
2. `pass` is the user's password, hashed with SHA-512.

Response:

```
{
	"success": {success},
	"message": "{message}",
	"position": {position}
}
```

1. `success` is a boolean indicating whether the operation was successful.
2. `message` is submitted to indicate an error, e.g. "Username already exists".
3. `position` is submitted if the operation was successful, but a reg_token
wasn't submitted, and represents the user's position in the queue of people
waiting for access to Lavaboom.

### Login

```
POST /login

{
	"username": "{user}",
	"password": {pass}
}
```

1. `success` is a boolean.
2. `password` is the submitted password, hashed with SHA-512.

Response:

```
{
	"success": "{success}",
	"message": "{message}",
	"token": "{token}",
	"exp_date": "{exp_date}"
}
```

1. `success` is a boolean.
2. `message` is an error message.
3. `token` is a UUID that acts as a session token.
4. `exp_date` is the date when the session token expires. This will be set pretty aggressively, due to the nature of our product.

Issues:

1. If the user is on mobile, the `exp_date` will be set much more permissively. A mobile login will be detected via HTTP headers.


### Logout

```
DELETE /logout
Auth: {token}
```

Response:

```
{
	"success": "{success}",
	"message": "{message}"
}
```
