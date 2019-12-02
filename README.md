# songjiang
[![Build Status](https://cloud.drone.io/api/badges/storezhang/songjiang/status.svg)](https://cloud.drone.io/storezhang/songjiang)
[![](https://images.microbadger.com/badges/image/storezhang/songjiang.svg)](https://microbadger.com/images/storezhang/songjiang "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/storezhang/songjiang.svg)](https://microbadger.com/images/storezhang/songjiang "Get your own version badge on microbadger.com")
[![](https://images.microbadger.com/badges/commit/storezhang/songjiang.svg)](https://microbadger.com/images/storezhang/songjiang "Get your own commit badge on microbadger.com")
[![Go Report Card](https://goreportcard.com/badge/github.com/storezhang/songjiang)](https://goreportcard.com/report/github.com/storezhang/songjiang)

# 自动签到，特点如下：
- 以Docker运行
- 极低的内存占用
- 极低的CPU消耗


# 为什么叫宋江
宋江，小说《水浒传》里面的梁山好汉的领袖，号及时雨，专干送钱的营生，给需要钱的人送去钱财，最后把梁山都送出去了。这个特性很像自动签到，及时补充账号中的钱财或者积分。


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
  --name=songjiang \
  storezhang/songjiang
```


# 配置
配置文件名为songjiang.yml，有如下配置项（**示例所配置的值为默认值**）
```
songjiang:
  debug: true # 是否开户Debug模式
  logLevel: debug # 日志级别

aliyun: # 阿里云的配置
  appKey: ${ALIYUN_APPKEY} # 阿里云的AppKey
  secret: ${ALIYUN_SECRET} # 阿里云的Secret

domains: # 域名配置
  - type: aliyun # 类型
    name: imyserver.com # 主域名
    subDomains: test # 子域名，以英文逗号,分隔
    subDomainPrefix: prefix # 子域名前缀，会加在所有的子域名前，如prefox.test
    subDomainStaff: staff # 子域名后缀，会加在所有子域名后，如test.staff
    value: ${VALUE} # 值，如果dnsType类型为CNAE需要填写此值
    dnsTypes: A # 域名类型，支持A,AAAA,CNAME等，以英文逗号,分隔
    redo: 1m # 执行间隔，支持1s、1m、1m1s等
  - type: aliyun # 类型
    name: imyserver.com # 主域名
    subDomains: test11,test22 # 子域名，以英文逗号,分隔
    value: test.imyserver.com # 值，如果dnsType类型为CNAE需要填写此值
    dnsType: CNAME A # 域名类型，支持A,AAAA,CNAME等，以英文逗号,分隔
    redo: 1h # 执行间隔，支持1s、1m、1m1s等
```
