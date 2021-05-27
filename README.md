- 初次需要 `token`，填写 `InitDB` 字段:
  ```go
  "clientID":           "",
	"clientSecretOnline": "",
	"refreshTokenOnline": "",
  ```

- 运行 `go run .`
- 打包 `go build .`，在windows打包到linux需要更改环境设置 `set GOOS=linux`
- linux后台运行 `nohup ./xxx >/dev/null 2>log &`