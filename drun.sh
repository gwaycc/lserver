# /bin/sh

# debug run
sudo docker run -it --rm \
    -p 11301:11301 -p 11302:11302 \
    -v $PRJ_ROOT/etc:/app/etc \
    -v $PRJ_ROOT/var/log:/app/var/log \
    -v $PRJ_ROOT/var/log:/app/var/log \
    -v $PRJ_ROOT/var/beanstalkd:/app/var/beanstalkd \
    $PRJ_NAME
