# ddns
[![Build Status](https://cloud.drone.io/api/badges/storezhang/ddns/status.svg)](https://cloud.drone.io/storezhang/ddns)
[![](https://images.microbadger.com/badges/image/storezhang/ddns.svg)](https://microbadger.com/images/storezhang/ddns "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/storezhang/ddns.svg)](https://microbadger.com/images/storezhang/ddns "Get your own version badge on microbadger.com")
[![](https://images.microbadger.com/badges/commit/storezhang/ddns.svg)](https://microbadger.com/images/storezhang/ddns "Get your own commit badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/storezhang/ddns)](https://goreportcard.com/report/github.com/storezhang/ddns)

# 自动签到，特点如下：
- 以Docker运行
- 极低的内存占用
- 极低的CPU消耗


# 如何使用
有丙种方法可以运行：
- **Docker（建议使用此方法）**
- ~~直接下载可执行程序（不建议）~~

直接使用命令行执行
```
sudo docker run \
  --volume=${YOUR_CONF_DIR}:/ddns \
  --restart=always \
  --detach=true \
  --name=ddns \
  storezhang/ddns
```


# 配置
application.yml或者application.toml，有如下配置项（**示例所配置的值为默认值**）
```
ddns:
  debug: false
  chans:
    - key: ${KEY}

aliyun:
  appKey: ${APPKEY}
  secret: ${SECRET}

domains:
  - type: aliyun
    redo: 1m
    name: xxx.com
    subDomains: yyy
    dnsTypes: A
  - type: aliyun
    redo: 1h
    name: xxx.com
    subDomains: aaa,bbb,ccc
    subDomainStaff: zzz
    value: yyy.xxx.com
    dnsTypes: CNAME
```
