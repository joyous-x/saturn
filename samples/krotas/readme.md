

# TODO

- [x] static
- [x] ip地域
    + [x] real ip for http.Request
- [ ] kv
- 登录
    - [ ] /v1/wx/miniapp/user/login
        D:\NRS\github\saturn\samples\krotas\common\rand.go
        commonReq,commonResp
    - [ ] /v1/wx/miniapp/user/update
    - [ ] /v1/wx/miniapp/access_token
    - [ ] /v1/user/login
        ```
        curl -XPOST 'http://localhost:8000/v1/user/login' -d '{}'
        ```
- 支付



## Note
- GO111MODULE=off 时，build项目：
    + 项目的 main.go 需要放在 src 下的第二级目录, 就是说: src/main.go 不行, src/ceggs/main.go 即可
