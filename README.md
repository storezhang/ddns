# ddns
[![编译状态](https://github.ruijc.com:20443/api/badges/storezhang/ddns/status.svg)](https://github.ruijc.com:20443/storezhang/ddns)
[![Golang质量](https://goreportcard.com/badge/github.com/storezhang/ddns)](https://goreportcard.com/report/github.com/storezhang/ddns)
![版本](https://img.shields.io/github/go-mod/go-version/storezhang/ddns)
![Docker镜像版本](https://img.shields.io/docker/v/storezhang/ddns)
![仓库大小](https://img.shields.io/github/repo-size/storezhang/ddns)
![Docker镜像大小](https://img.shields.io/docker/image-size/storezhang/ddns)
![最后提交](https://img.shields.io/github/last-commit/storezhang/ddns)
![授权协议](https://img.shields.io/github/license/storezhang/ddns)
![星星个数](https://img.shields.io/github/stars/storezhang/ddns?style=social)

# 动态域名解析，特点如下：
- 原生`Docker`支持
- 极低的内存占用
- 极低的`CPU`消耗，实测几十个域名解析也消耗不到`0.1%`的性能
- 多域名厂商支持
  - 阿里云
  - 腾讯云
- 配置文件监控
- 多域名支持


# 如何使用

直接使用命令行执行
```shell
TAG="ccr.ccs.tencentyun.com/storezhang/ddns" && NAME="Ddns" && sudo docker pull ${TAG} && sudo docker stop ${NAME} ; sudo docker rm --force --volumes ${NAME} ; sudo docker run \
  \
  \
  \
  --volume=/主机目录:/config \
  --volume=/etc/localtime:/etc/localtime \
  \
  \
  \
  --env=UID=$(id -u 用户名) \
  --env=GID=$(id -g 用户名) \
  \
  \
  \
  --restart=always \
  --detach=true \
  --name=${NAME} \
  ${TAG} \
  \
  \
  \
&& sudo docker logs -f ${NAME}
```

# 配置

默认的配置文件如下
```yaml
secrets:
  - ak: ${ALIYUN_AK}
    sk: ${ALIYUN_SK}
    # 目前只支持阿里云
    type: aliyun
    # 任意字符，如果有多个，后续和解析绑定
    # 可以不配置，有默认值，那样所有解析都使用本授权
    label: test-label

resolves:
  - name: ruijc.com
    # 绑定授权，可以不配置
    label: test-label
    types:
      - CNAME
    value: storezhang.ruijc.com
    subdomains:
      - test
  - name: ruijc.com
    types:
      - A
    subdomains:
      - test
```

配置文件可以使用`${ENV}`来加载环境变量
