FROM alpine:3.5

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add tzdata ca-certificates && cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /var/cache/apk/*
ADD ./ /data/apps/krotas/

EXPOSE 8000
EXPOSE 8001

CMD [ "/data/apps/krotas/krotas", "-config", "/data/apps/krotas/config"]
