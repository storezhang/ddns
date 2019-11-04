FROM storezhang/chromium:73.0.3683

ENV YANGJIAN_HOME /yangjian-data

ADD yangjian /opt/yangjian/

RUN set -x \
    && apk update \
    && apk --no-cache add bash \
    && apk --no-cache add curl \
    && apk --no-cache add openssl \
    && apk --no-cache add libidn \
    && chmod +x /opt/yangjian/yangjian \
    && rm -rf /var/cache/apk/*

VOLUME ${YANGJIAN_HOME}
WORKDIR ${YANGJIAN_HOME}

ENTRYPOINT ["/opt/yangjian/yangjian"]
