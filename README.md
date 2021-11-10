1. mysql 数据库表 `global` 初始化：
   | gid | BASE_URL | CLIENT_ID | CLIENT_SECRET | API_KEY | SECRET_KEY | API_URL TOKEN_URL |
   | :--:| :------: | :-------: | :-----------: | :-----: | :--------: | :---------------: |
   | 1 | https://graph.microsoft.com/v1.0/me | `microsoft azure` 后台添加程序获得 | `microsoft azure` 后台添加程序获得 | `百度智能云`注册应用获取 | `百度智能云`注册应用获取 | `百度智能云`注册应用获取 |
2. mysql 数据库表 `user` 初始化：
   | uid | name | email | photo | refresh | access | uptime |
   | :-: | :--: | :---: | :---: | :-----: | :----: | :----: |
   | 1 | `microsoft azure`获取 | `microsoft azure`获取 | `microsoft azure`获取(longText) | `rclone`获取(longText) | `rclone`获取(longText) | bigint |
3. mysql 数据库表 `essay` 初始化：
   | aid | essayId | title | size | note | content | addtime | uptime |
   | :-: | :-----: | :---: | :--: | :--: | :-----: | :-----: | :----: |
   | 自增 | varchar | varchar | varchar | longText | longText | bigint | bigint |
4. mysql 数据库表 `one` 初始化：
   | oid | md |
   | :-: | :-: |
   | 自增 | longText |
5. 填写完善`config.DB`，连接数据库。初始化`one`数据，可以运行`global.SetOne`函数。

=================
- 运行 `go run main.go`  

- 打包 `go build .`，在windows打包到linux需要更改环境设置 `go env -w GOOS=linux`  

- linux后台运行 `nohup ./xxx >/dev/null 2>log &`