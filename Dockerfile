FROM storezhang/chromium


MAINTAINER storezhang "storezhang@gmail.com"
LABEL architecture="AMD64/x86_64" version="latest" build="2019-12-19"
LABEL Description="基于Alpine的自动签到镜像，支持Hao4K这类主流网站，也提供了ServerChan推送。"


ENV USERNAME songjiang
ENV ROOT_DIR /songjiang
ENV UID 1000
ENV GID 1000


WORKDIR ${ROOT_DIR}
VOLUME ${ROOT_DIR}


ADD songjiang /opt
COPY docker /


RUN set -ex \
    \
    && addgroup -g ${GID} -S ${USERNAME} \
    && adduser -u ${UID} -g ${GID} -S ${USERNAME} \
    \
    && apk update \
    \
    && mkdir -p /conf \
    && chmod +x /usr/bin/entrypoint \
    && chmod +x /etc/s6/.s6-svscan/* \
    && chmod +x /etc/s6/songjiang/* \
    \
    && apk --no-cache add bash s6 \
    && rm -rf /var/cache/apk/*


ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/bin/s6-svscan", "/etc/s6"]
