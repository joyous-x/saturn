- /wx/miniapp/user/login
```
curl -XPOST "http://localhost:8000/wx/miniapp/user/login" -d '{"common":{"uid":"test-uid", "appid":"test-appid"}, "jscode": "", "inviter":""}'
```

- /wx/miniapp/user/update
```
curl -XPOST "http://localhost:8000/wx/miniapp/user/update" -d '{"common":{"uid":"test-uid", "appid":"test-appid"}, "encryptedData": "", "iv":""}'
```

- /wx/miniapp/user/access_token
```
curl -XPOST "http://localhost:8000/wx/miniapp/user/access_token" -d '{"common":{"uid":"test-uid", "appid":"test-appid"}, "appid": ""}'
```

- /wx/miniapp/user/event_message
```
curl -XPOST "http://localhost:8000/wx/miniapp/user/event_message" -d '{"common":{"uid":"test-uid", "appid":"test-appid"}, "appid": ""}'

curl -XGET "http://localhost:8000/wx/miniapp/user/event_message" -d '{"common":{"uid":"test-uid", "appid":"test-appid"}, "appid": ""}'
```