- 初次需要 `token`，填写 `InitDB` 字段:
  ```go
    "clientID":           "",
  "clientSecretOnline": "",
  "refreshTokenOnline": "",
  ```
- `refreshToken` 可以通过 `rclone` 获得
- `clientID` 和 `clientSecret`，通过 `microsoft azure` 后台添加程序获得


=================
- 运行 `go run .`
- 打包 `go build .`，在windows打包到linux需要更改环境设置 `set GOOS=linux`
- linux后台运行 `nohup ./xxx >/dev/null 2>log &`