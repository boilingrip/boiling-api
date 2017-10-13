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
GET /artists/autocomplete/{s}
GET /artist/autocomplete_tags/{s}

GET /release_groups/{id}

GET /formats
GET /leech_types
GET /media
GET /release_group_types
GET /release_properties
GET /release_roles
GET /privileges
```

### The `/login` and `/signup` Endpoints

See the Calling the API section above.

### The `GET /artists/{id}` Endpoint

The `/artists/{id}` endpoint returns the artist with the given ID.
This endpoint requires the `get_artist` privilege.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/artists/1'
```

Response:
```json
{"status":"success","data":{"artist":{"id":1,"name":"Led Zeppelin","added":"2017-10-13T21:41:31.411901Z","added_by":{"id":1,"username":"test"},"bio":"Some American Band","tags":["rock","70s","80s","usa"]}}}
```

### The `GET /artists/autocomplete/{s}` Endpoint

The `/artists/autocomplete/{s}` endpoint returns a list of artists for auto-completion of an artists name or alias.
This endpoint requires the `get_artist` privilege.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/artists/autocomplete/Led'
```

Response:
```json
{"status":"success","data":{"artists":[{"id":1,"name":"Led Zeppelin","added":"2017-10-13T21:41:31.411901Z","added_by":{"id":1,"username":"test"},"bio":"Some American Band","tags":["rock","70s","80s","usa"]}]}}
```

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/artists/autocomplete/e'
```

Response:
```json
{"status":"success","data":{"artists":[{"id":1,"name":"Led Zeppelin","added":"2017-10-13T21:41:31.411901Z","added_by":{"id":1,"username":"test"},"bio":"Some Canadian producer of electronic music","tags":["rock","70s","80s","usa"]},{"id":2,"name":"deadmau5","aliases":[{"alias":"testpilot","added":"2017-10-13T21:41:31.473056Z","added_by":{"id":1,"username":"test"}}],"added":"2017-10-13T21:41:31.411901Z","added_by":{"id":1,"username":"test"},"bio":"Some Canadian producer of electronic music","tags":["2000s","2010s","canada","edm","techno"]}]}}
```

### The `GET /artists/autocomplete_tags/{s}` Endpoint

The `/artists/autocomplete_tags/{s}` endpoint returns a list of artist tags auto-completed from `s`.
This endpoint requires the `get_artist` privilege.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/artists/autocomplete_tags/garbage'
```

Response:
```json
{"status":"success","data":{"tags":null}}
```

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/artists/autocomplete_tags/c'
```

Response:
```json
{"status":"success","data":{"tags":["rock","canada","techno"]}}
```

### The `GET /release_groups/{id}` Endpoint

The `/release_groups/{id}` endpoint returns the release group with the given ID.
This endpoint requires the `get_release_group` privilege.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/release_groups/1'
```

Response:
```json
{"status":"success","data":{"release_group":{"id":1,"name":"4x4=12","artists":[{"role":"Main","artist":{"id":2,"name":"deadmau5"}}],"release_date":"2010-12-03T00:00:00Z","added":"2017-10-13T21:41:31.500883Z","added_by":{"id":1,"username":"test"},"type":"Album","tags":["edm","techno","electronic"]}}}
```

### The `/formats` Endpoint

The `/formats` endpoint returns a list of all possible formats.
Each format is a concatenation of the actual format (for example `MP3/V0`) and its encoding (for example `Lossy`).
The two values can be split at the dollar sign, for now.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/formats'
```

Response:
```json
{"status":"success","data":{"formats":["FLAC$Lossless","FLAC/24bit$Lossless","MP3/320$Lossy","MP3/V0$Lossy","MP3/V2$Lossy"]}}
```

### The `/leech_types` Endpoint

The `/leech_types` endpoint returns a list of all possible leech types.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/leech_types'
```

Response:
```json
{"status":"success","data":{"leech_types":["DoubleDown","DoubleUp","Freeleech","Neutral","Normal"]}}
```

### The `/media` Endpoint

The `/media` endpoint returns a list of all possible media.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/media'
```

Response:
```json
{"status":"success","data":{"media":["Blu-Ray","CD","Cassette","DAT","DVD","SACD","Vinyl","WEB"]}}
```

### The `/release_group_types` Endpoint

The `/release_group_types` endpoint returns a list of all possible release_group_types.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/release_group_types'
```

Response:
```json
{"status":"success","data":{"release_group_types":["Album","Bootleg","Compilation","EP","Live album","Mixtape","Single","Soundtrack","Unknown"]}}
```

### The `/release_properties` Endpoint

The `/release_properties` endpoint returns a list of all possible release properties.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/release_properties'
```

Response:
```json
{"status":"success","data":{"release_properties":["CassetteApproved","LossyMasterApproved","LossyWebApproved"]}}
```

### The `/release_roles` Endpoint

The `/release_roles` endpoint returns a list of all possible release roles.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/release_roles'
```

Response:
```json
{"status":"success","data":{"release_roles":["Composer","Conductor","Guest","Main","Producer","Remixer"]}}
```

### The `/privileges` Endpoint

The `/privileges` endpoint returns a list of all possible privileges.
No privileges a re required for this endpoint.

Request:
```bash
curl -X GET -H 'X-User-Token: <elided>' 'http://localhost:8080/privileges'
```

Response:
```json
{"status":"success","data":{"privileges":["delete_blog","delete_blog_not_owner","get_artist","get_blogs","get_release_group","post_blog","post_blog_override_author","post_blog_override_posted_at","update_blog","update_blog_not_owner","update_blog_override_author","update_blog_override_posted_at"]}}
```
