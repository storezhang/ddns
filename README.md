# yangjian（杨戬）
[![Build Status](https://cloud.drone.io/api/badges/storezhang/ddns/status.svg)](https://cloud.drone.io/storezhang/ddns)

支持DDNS（动态域名解析），特点如下：
- 支持多域名
- 支持多提供商
- 支持Docker运行
- 极低的内存占用
- 极低的CPU占用


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
