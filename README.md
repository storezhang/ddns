# yangjian（杨戬）
[![Build Status](https://drone.storezhang.imyserver.com:20443/api/badges/imynas/yangjian/status.svg)](https://drone.storezhang.imyserver.com:20443/imynas/yangjian)

基于ACME协议的NAS自动HTTPS更新，使用群晖自带的管理后台进行证书更新，所以不会有任何问题，支持群晖所有套件的一件证书更新，包括：
- LDAP服务器
- Drive
- FTPS
- WebDAV服务器
- 其它所有群晖支持的需要HTTPS套件


# 支持的NAS系统
- 群晖（Synology）


# 为什么要取名叫杨戬
在现行的文化中，中国人对自己的文化不自信（参看头条上一说到外国男人和中国女人下面就一群人留言什么长短什么什么的），所以命名全部以中国
名人来命名，增强文化上的自信。且在中国神话中，杨戬是最著名的**督粮官**，而HTTPS对于NAS来说，就等同于一只军队的粮草，是非常重要的！


# 如何使用
有丙种方法可以运行：
- **Docker（建议使用此方法）**
- ~~直接下载可执行程序（不建议）~~

直接使用命令行执行
```
docker run \
  --volume=/your/config/path:/yangjian-data \
  --name=yangjian \
  storezhang/yangjian
```


# 配置
配置文件名为config.yml，有如下配置项（**示例所配置的值为默认值**）
```
db:
  dsn: "root:@tcp(127.0.0.1:3306)/niulang?charset=utf8&parseTime=True&loc=Local"  # 数据库的dsn
  backup: false # 升级前是否备份数据库

niulang:
  port: 8000  #http server的端口
  debug: false # 是否开启调试模式，在Docker下一定要为false，不然无法运行
  browserWidth: 2560 # 窗口宽度
  browserHeight: 1600 # 窗口调试
  execDurationMonth: 24h # 运行间隔
  browserTimeout: 1h # 浏览器超时时间
  timeFormat: 2006-01-02 15:04:05 # 参看Golang的时间格式配置
  delayExit: 15s # 结束后多久关闭浏览器

ssls:
  - domain: nas.imyserver.com # 你要设置的NAS的域名
    url: https://nas.imyserver.com # 你的群晖登录地址，可以是内网地址
    username: xxxx # 群晖的用户名（确保此用户能操作证书）
    password: xxxx # 群晖的密码
    type: synology # NAS的类型，支持类型synology（群晖）、unas（UNAS）
    lang: English # 群晖语言（请在你的群晖系统区域设置里查看有哪些语言，如English、简体中文等）
    acme: # ACME配置，具体配置参看ACME的配置
        dns: dns_ali # 域名提供商的类型
        dnsSleep: 120 # DNS配置间隔时间，用的是ACME的DNS域名验证方式，此时间为DNS生效时间
        aliKey: xxxx # 阿里的Key
        aliSecret: xxxx # 阿里的Secret

```

关于ACME的配置，可以参看[ACME官网](https://github.com/Neilpang/acme.sh)来做配置配置，配置项可以使用驼峰形式，也可以使用ACME文件里面的配置，举个粟子：
- Ali_Key可以写成aliKey、Ali_Key
- Ali_Secret可以写成aliSecret、Ali_Secret
- DOMAIN可以写成DOMAIN、domain

**程序会对上面的配置做出自动配置**

# API

API 使用JWT认证方式

## 登录API
URL:
```bash
http://127.0.0.1:8000/login
```
返回：
```json
{
    "code": 200,
    "msg": "ok",
    "data": {
        "expired": "2019-08-13T18:10:17.928207+08:00",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjU2OTEwMTcsImlkIjoidGVzdDUiLCJvcmlnX2lhdCI6MTU2NTY4NzQxN30.Z4OGgPy5TD86plQkdKm7tX2H_fUF5ZtfMzdz8o-QAPM"
    }
}
```

## 用户API, 暂时支持 GET, POST, DELETE
URL:
```bash
http://127.0.0.1:8000/user/users/8
```
数据结构
```json
{
    "id": 1,
    "username": "admin",
    "password": "******",
    "lastName": "",
    "firstName": ""
}
```