
-- 后台管理用户表
CREATE TABLE mail_server
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

