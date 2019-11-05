FROM storezhang/alpine


MAINTAINER storezhang "storezhang@gmail.com"
LABEL architecture="AMD64/x86_64" version="latest" build="2019-11-05"
LABEL Description="基于Alpine的DDNS功能镜像，支持阿里云、DNSPod等主流DDNS厂商。"


ENV USERNAME ddns
ENV UID 1000
ENV GID 1000


WORKDIR /
VOLUME ["/data"]
VOLUME ["/conf"]


ADD ddns /opt


RUN set -ex \
    \
    && addgroup -g ${GID} -S ${USERNAME} \
    && adduser -u ${UID} -g ${GID} -S ${USERNAME} \
    \
    && apk update \
    \
    && mkdir -p /conf \
    && mkdir -p /data \
    \
    && apk --no-cache add bash s6 \
    && rm -rf /var/cache/apk/*


COPY root /


ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/bin/s6-svscan", "/etc/s6"]

