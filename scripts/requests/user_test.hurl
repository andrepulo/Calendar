POST http://{{host}}/api/v1/signin
Authorization: Bearer test

HTTP 401

POST http://{{host}}/api/v1/verify
Authorization: Bearer test

HTTP 403

POST http://{{host}}/api/v1/signup
[BasicAuth]
test: test

HTTP 200
[Captures]
token: header "Authorization"

POST http://{{host}}/api/v1/signin
[BasicAuth]
test: test

HTTP 200
[Captures]
token: header "Authorization"

POST http://{{host}}/api/v1/verify
Authorization: {{token}}

HTTP 200
[Asserts]
jsonpath "$.id" exists
jsonpath "$.login"  == "test"
[Captures]
user_id: jsonpath "$['id']"


PATCH http://{{host}}/api/v1/users
Authorization: {{token}}
{
   "password": "new_test"
}
HTTP 200

POST http://{{host}}/api/v1/signin
[BasicAuth]
test: test

HTTP 401

POST http://{{host}}/api/v1/signin
[BasicAuth]
test: new_test

HTTP 200
