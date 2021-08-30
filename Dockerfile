FROM golang:1.16-alpine as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

RUN mkdir -p /tmp/go-plugins/ && \
    sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories

COPY ./kong-rbac-plugin /tmp/go-plugins/

RUN cd /tmp/go-plugins/ && \
    apk add make libc-dev gcc && \
    make all


FROM kong:2.5.0-alpine

COPY --from=builder  /tmp/go-plugins/go-pluginserver /usr/local/bin/go-pluginserver
COPY --from=builder  /tmp/go-plugins/rbac /usr/local/bin/rbac

RUN /usr/local/bin/go-pluginserver -version && /usr/local/bin/rbac -dump
