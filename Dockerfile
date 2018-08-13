FROM alpine:3.7
MAINTAINER SHU <free1139@163.com>

## Add and Beanstalkd
RUN apk add --update ca-certificates beanstalkd redis
RUN mkdir -p /app/var/beanstalkd/
RUN mkdir -p /app/var/redis/
RUN mkdir -p /app/var/log/
COPY $PRJ_ROOT/bin/docker/bsa /usr/bin/
COPY $PRJ_ROOT/bin/docker/supd /usr/bin/
COPY $PRJ_ROOT/publish/lserver/ /app

EXPOSE 11300
EXPOSE 11302

CMD ["supd", "-c","/app/etc/supervisord.conf", "-n"]

