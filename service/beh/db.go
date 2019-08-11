package main

import (
	"fmt"
	"sync"
	"time"

	"lserver/module/db"
	"lserver/module/gouuid"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/behavior"
)

const (
	putBehaviorSql = `
    INSERT INTO
    %s
    (
       event_time, event_key, from_ip, req_method, req_body,
       resp_code, resp_body, use_nsec, uuid
    )VALUES(
            ?, ?, ?, ?, ?,
            ?, ?, ?, ?
    )
`
	createBehaviorTb = `
CREATE TABLE IF NOT EXISTS %s
(
    -- 产生事件时的时间
    -- 事件产生器必须生成事件时间
    event_time TIMESTAMP NOT NULL,
    -- 读取的key值，如果有用户ID填写用户ID, 视需要填写
    event_key VARCHAR(128) NOT NULL,

	-- 请求来源ip
	from_ip VARCHAR(128),
    -- 请求方法
    req_method VARCHAR(40),
    -- 请求数据
    req_body BLOB,
    -- 响应的状态码
    resp_code VARCHAR(10),
    -- 响应的内容
    resp_body BLOB,
    -- 事件用时
    -- 如果事件结束时间， 用时为nsec
    use_nsec BIGINT, 

    -- 唯一值, 用于避免主键重复, 此表实际上只需索引，不需主键
    uuid VARCHAR(128),
    -- 主键
    PRIMARY KEY(event_time, uuid),
	KEY(event_key, event_time, resp_code)
) ENGINE=InnoDB;
    `
)

var (
	behaviorTbName = "beh_"
	behaviorTbTime = "" // 201510
	behaviorTbLock = sync.Mutex{}
)

func getBehaviorTbName(currentTime time.Time) (string, error) {
	behaviorTbLock.Lock()
	defer behaviorTbLock.Unlock()
	timefmt := currentTime.Format("200601")
	// 时间相等，说明已经检查过分表名
	if behaviorTbTime == timefmt {
		return behaviorTbName + behaviorTbTime, nil
	}
	newTbName := behaviorTbName + timefmt
	mdb := db.GetCache("master")
	// 创表
	if _, err := mdb.Exec(fmt.Sprintf(createBehaviorTb, newTbName)); err != nil {
		return "", errors.As(err)
	}
	// 创建成功后再刷新缓存的时间
	behaviorTbTime = timefmt
	return newTbName, nil
}

// 生成记录
func insertBehavior(l *behavior.Event) error {
	// 获取分表名
	tbName, err := getBehaviorTbName(l.EventTime)
	if err != nil {
		return errors.As(err, l)
	}
	methodLen := len(l.ReqMethod)
	if methodLen > 40 {
		l.ReqMethod = l.ReqMethod[:40]
	}
	mdb := db.GetCache("master")
	_, err = mdb.Exec(fmt.Sprintf(putBehaviorSql, tbName),
		l.EventTime,
		l.EventKey,
		l.FromIp,
		l.ReqMethod,
		l.ReqBody,
		l.RespCode,
		l.RespBody,
		l.UseTime,
		gouuid.New(),
	)
	if err != nil {
		behaviorTbTime = "" // 重新检查表数据是否创建了
		return errors.As(err, l)
	}
	return nil
}
