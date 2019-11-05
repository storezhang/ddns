# ddns
[![Build Status](https://cloud.drone.io/api/badges/storezhang/ddns/status.svg)](https://cloud.drone.io/storezhang/ddns)
[![](https://images.microbadger.com/badges/image/storezhang/ddns.svg)](https://microbadger.com/images/storezhang/ddns "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/storezhang/ddns.svg)](https://microbadger.com/images/storezhang/ddns "Get your own version badge on microbadger.com")
[![](https://images.microbadger.com/badges/commit/storezhang/ddns.svg)](https://microbadger.com/images/storezhang/ddns "Get your own commit badge on microbadger.com")

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
sudo docker run \
  --net=host \
  --volume=${YOUR_DDNS_DIR}:/conf \
  --restart=always \
  --detach=true \
  --name=ddns \
  storezhang/ddns
```


# 配置
配置文件名为ddns.yml，有如下配置项（**示例所配置的值为默认值**）
```
ddns:
  debug: true # 是否开户Debug模式
  logLevel: debug # 日志级别
  redo: 1m # 执行间隔，支持1s、1m、1m1s等

aliyun: # 阿里云的配置
  appKey: ${ALIYUN_APPKEY} # 阿里云的AppKey
  secret: ${ALIYUN_SECRET} # 阿里云的Secret

domains: # 域名配置
  - type: aliyun # 类型
    name: imyserver.com # 主域名
    subDomains: test # 子域名，以英文逗号,分隔
    dnsTypes: A # 域名类型，支持A,AAAA,CNAME等，以英文逗号,分隔
```
