# 日志管理服务项目

客户端向beanstalkd的log.gway.cc队列发送协议日志，以便可以统一监控与管理多个项目的日志情况。

管理界面：http://host:11302

TODO:重构后台界面

# 依赖安装
```text
# debian环境
sudo aptitude install mysql-server
sudo aptitude install redis
sudo aptitude install beanstalkd
```

# 数据库初始化
```
cd $PRJ_ROOT/doc/sql
# 修改init.sh的host、用户名、passwd
./init.sh

# 修改$PRJ_ROOT/etc/db.cfg下username与passwd与实际相符以便代码中使用
```

# sup部署

## 加载环境变量
```text
source env.bash # 加载项目环境变量
```

## 编译与部署
```text
sup build all
sup install all
sup status
```

# docker部署

## 安装docker
```text
# debian环境
sudo aptitude install docker
```

## 加载环境变量
```text
source env.bash # 加载项目环境变量
```

## 编译docker镜像
```text
cd $PRJ_ROOT
./dbuild.sh
```

## 启动部署
### 启动beanstalkd镜像(若本地已有，不需再启动)
```text
# 使用自带配置
[sudo] docker run -d --restart=always -p 11301:11301 --name=$PRJ_NAME.beanstalkd $PRJ_NAME beanstalkd -p 11301 -b /app/var/beanstalkd

# 使用外部配置
[sudo] docker run -d --restart=always -p 11301:11301 -v $PRJ_ROOT/etc:/app/etc -v $PRJ_ROOT/var:/app/var --name=$PRJ_NAME.beanstalkd $PRJ_NAME beanstalkd -p 11301 -b /app/var/beansatalkd
```

### 启动redis服务(若本地已有，不需再启动)
```text
# 使用自带配置
[sudo] docker run -d --restart=always -p 4932:4932 --name=$PRJ_NAME.redis $PRJ_NAME redis-server

# 使用外部配置
[sudo] docker run -d --restart=always -p 4932:4932 -v $PRJ_ROOT/etc:/app/etc -v $PRJ_ROOT/var:/app/var --name=$PRJ_NAME.redis $PRJ_NAME redis-server
```

### 启动服务程序
```text
# 使用自带配置
[sudo] docker run -d --restart=always --name=$PRJ_NAME.service.log $PRJ_NAME /app/src/service/log
[sudo] docker run -d --restart=always --name=$PRJ_NAME.applet.cms $PRJ_NAME /app/src/applet/cms

# 使用外部配置
[sudo] docker run -d --restart=always -v $PRJ_ROOT/etc:/app/etc -v $PRJ_ROOT/var:/app/var --name=$PRJ_NAME.service.log $PRJ_NAME /app/src/service/log
[sudo] docker run -d --restart=always -v $PRJ_ROOT/etc:/app/etc -v $PRJ_ROOT/var:/app/var --name=$PRJ_NAME.applet.cms $PRJ_NAME /app/src/applet/cms
```

# 附录
[协议设计说明](https://github.com/gwaycc/lserver/tree/master/doc)
[客户端实现](https://github.com/gwaylib/log) 
