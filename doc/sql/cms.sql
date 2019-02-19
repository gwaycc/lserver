
-- 创建数据库
CREATE SCHEMA `lserver` DEFAULT CHARACTER SET utf8mb4 ;
USE lserver;

-- 后台管理用户表
CREATE TABLE cms_user
(
    -- 用户名
    username VARCHAR(32) NOT NULL,
    -- 创建时间
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 用户密码
    `passwd` VARCHAR(128) NOT NULL,
    -- 用户昵称
    nickname VARCHAR(128) NOT NULL DEFAULT '',
    -- 用户组，０为管理员，１为非管理员
    -- 根据使用情况，暂不实现组权限功能
    group_id INT NOT NULL DEFAULT 1,
    -- 1，可用，2, 禁用。
    status INT NOT NULL DEFAULT 1,
    -- 主键
    PRIMARY KEY(username)
);
-- 默认密码LogAdmin123
INSERT INTO cms_user(username,`passwd`, nickname, group_id)VALUES('admin','$2a$10$4B9dd5YdxEHoN/I9LKGJhuheh7pqM2smhO0vXUEWpVeJWNCI0TL0O','管理员',0);

-- 后台操作记录
CREATE TABLE cms_log
(
    -- 操作时间
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- 操作的用户
    username VARCHAR(32) NOT NULL,
    -- 操作类别
    kind VARCHAR(128) NOT NULL,
    -- 输入的参数
    args TEXT NOT NULL,
    -- 备注
    memo TEXT NOT NULL,
    -- 主键
    PRIMARY KEY(created_at, username, kind)
);

-- 后台菜单表
CREATE TABLE cms_menu
(
    -- 菜单键值
    -- 一级菜单，user
    -- 二级菜单，user.create, user.get等
    id VARCHAR(32) NOT NULL,
    -- 菜单名称
    name VARCHAR(128) NOT NULL,
    -- 主键
    PRIMARY KEY(id)
);
INSERT INTO cms_menu(id,name)VALUES('*.*.*.*.*','所有权限');
INSERT INTO cms_menu(id,name)VALUES('app.cms.info','后台.用户查询');
INSERT INTO cms_menu(id,name)VALUES('app.cms.create','后台.用户创建');
INSERT INTO cms_menu(id,name)VALUES('app.cms.pwd','后台.重置密码');
INSERT INTO cms_menu(id,name)VALUES('app.cms.delete','后台.用户删除');
INSERT INTO cms_menu(id,name)VALUES('app.cms.log','后台.操作日志');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv','后台.权限.查询');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.bind','后台.权限.快速设定');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.on','后台.权限.开通');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.off','后台.权限.关闭');
INSERT INTO cms_menu(id,name)VALUES('app.cms.privtpl','后台.权限.模板.模询');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.tpl.list','后台.权限.模板.列表');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.tpl.on','后台.权限.模板.开通');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.tpl.off','后台.权限.模板.关闭');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.tpl.new','后台.权限.模板.新建');
INSERT INTO cms_menu(id,name)VALUES('app.cms.priv.tpl.delete','后台.权限.模板.删除');
INSERT INTO cms_menu(id,name)VALUES('app.dashboard','主页');
INSERT INTO cms_menu(id,name)VALUES('app.log.info','平台日志.日志查询');
INSERT INTO cms_menu(id,name)VALUES('app.log.alertor','平台日志.告警联系人');
INSERT INTO cms_menu(id,name)VALUES('app.log.alertor.add','平台日志.告警联系人.增加');
INSERT INTO cms_menu(id,name)VALUES('app.log.alertor.set','平台日志.告警联系人.修改');
INSERT INTO cms_menu(id,name)VALUES('app.log.alertor.del','平台日志.告警联系人.删除');
INSERT INTO cms_menu(id,name)VALUES('app.log.mail','平台日志.邮件设置');
INSERT INTO cms_menu(id,name)VALUES('app.log.mail.set','平台日志.邮件设置.修改');

-- 后台权限表,　直接更改数据库无效，需用户重登录才生效
-- 权限存在即可访问
CREATE TABLE cms_user_priv
(
    -- 用户名
    username VARCHAR(32) NOT NULL,
    -- 菜单名
    menu_id VARCHAR(32) NOT NULL,
    -- 主键
    PRIMARY KEY(username, menu_id)
);
INSERT INTO cms_user_priv(username,menu_id)VALUES("admin", "*.*.*.*.*");

-- 后台权限模板表,　直接更改数据库无效，需用户重登录才生效
-- 权限存在即可访问
CREATE TABLE cms_user_priv_tpl
(
    -- 模板名称
    tplname VARCHAR(32) NOT NULL,
    -- 开通的菜单名
    menu_id VARCHAR(32) NOT NULL,
    -- 主键
    PRIMARY KEY(tplname, menu_id)
); 
INSERT INTO cms_user_priv_tpl SELECT '管理员模板' AS tplname, id AS menu_id FROM cms_menu WHERE id <> '*.*.*.*.*';

-- 内存配置文件，每5分钟读取一次
CREATE TABLE lserver_cfg
(
    -- 配置文件名称
    cfgname VARCHAR(32) NOT NULL,
    -- 配置文件内容
    cfgdata BLOB NOT NULL,
    -- 主键
    PRIMARY KEY(cfgname)
); 

