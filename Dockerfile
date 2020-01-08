FROM storezhang/alpine


MAINTAINER storezhang "storezhang@gmail.com"
LABEL architecture="AMD64/x86_64" version="latest" build="2020-01-08"
LABEL Description="基于Alpine的DDNS镜像，支持阿里云、百度云、腾讯云、DNSPod等。"


ENV USERNAME ddns
ENV ROOT_DIR /ddns
ENV UID 1000
ENV GID 1000


WORKDIR ${ROOT_DIR}
VOLUME ${ROOT_DIR}


ADD ddns /opt
COPY docker /


RUN set -ex \
    \
    && addgroup -g ${GID} -S ${USERNAME} \
    && adduser -u ${UID} -g ${GID} -S ${USERNAME} \
    \
    && apk update \
    \
    && mkdir -p ${ROOT_DIR} \
    && chmod +x /usr/bin/entrypoint \
    && chmod +x /etc/s6/.s6-svscan/* \
    && chmod +x /etc/s6/ddns/* \
    && chmod +x /opt/ddns \
    \
    && rm -rf /var/cache/apk/*


ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/bin/s6-svscan", "/etc/s6"]
