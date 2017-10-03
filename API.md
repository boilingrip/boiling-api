# API Documentation [Draft :)]

This shall document the API.

## Calling the API

Every call to the API must be made with the header `X-User-Token=<user token>`.
The user token can be obtained by logging in with valid credentials.

TODO we might want to add a per-app key.

An example of how to log in:
```bash
$ curl -X POST -F 'username=test' -F 'password=test' http://localhost:8080/login
{"status":"success","data":{"user":{"id":1,"username":"test","email":"test@boiling.rip","bio":"","enabled":true,"can_login":true,"joined_at":"2000-01-01T00:00:00Z","last_login":"2000-01-01T00:00:00Z","last_access":"2000-01-01T00:00:00Z","uploaded":0,"downloaded":0},"token":"873c8b94de8251a63898c41451348a6d1fd436fb782db321d3b7d4d493a94a6589141179765103f7072ebbfc7bfbedc14836ca452d4c6bcf0f7a7de9ce83779d"}}
```

A failed login attempt:
```bash
$ curl -X POST -F 'username=test' -F 'password=testa' http://localhost:8080/login
{"status":"fail","message":"invalid password"}
```

An example of a successive call, after being logged in:
```bash
$ curl -X GET -H 'X-User-Token: 873c8b94de8251a63898c41451348a6d1fd436fb782db321d3b7d4d493a94a6589141179765103f7072ebbfc7bfbedc14836ca452d4c6bcf0f7a7de9ce83779d' 'http://localhost:8080/blogs?offset=0&limit=100'
{"status":"success","data":{"entries":[]}}
```

An example of a call with an invalid token:
```bash
$ curl -X GET -H 'X-User-Token: garbage' 'http://localhost:8080/blogs?offset=0&limit=100'
{"status":"fail","message":"invalid token"}
```

### Responses

Every response is a `Response` struct, containing at least the field `status`.
If `status==success`, the call was successful and an optional `data` field contains the response.
If `status==fail`, the call was unsuccessful due to the user's fault and a `message` field contains the description of that failure.
If `status==error`, the call was unsuccessful due to a server-side error and an optional `message` field contains a description of that error.

## Types

See the `api` package for now.

## Methods

Endpoints:

```
POST /login with form username=asdf password=asdf
POST /signup with form username=asdf password=asdf email=asdf

GET /blogs?limit=50&offset=0
maybe? GET /blogs/{id}
POST /blogs < Form (create)
POST /blogs/{id} < Form (update)
DELETE /blogs/{id}

GET /users (self)
GET /users/{id}
POST /users/{id} < Form (update)
POST /users < Form (create, as admin?)

GET /artists/{id}
```