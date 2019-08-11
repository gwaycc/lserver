package main

import (
	"fmt"
	"sync"
	"time"

	"lserver/module/db"

	"github.com/gwaylib/errors"
)

const (
	putLogSql = `
INSERT INTO
%s
(
    md5, platform, version, ip, date, 
    level, logger, msg
)VALUES(
        ?, ?, ?, ?, ?,
        ?, ?, ?
)
`
	createLogTb = `
CREATE TABLE IF NOT EXISTS %s
(
    md5 char(64),

	-- platform name
	platform char(32) NOT NULL,

	-- platform version
	version char(32) NOT NULL,

	-- platform server at
	ip char(64),

	-- log date time
	date timestamp NOT NULL,

	-- log level
	level INT NOT NULL,

	-- logger name
	logger char(64) NOT NULL, 

	-- log message
	msg BLOB NOT NULL, 

	PRIMARY KEY (md5),
	KEY (platform,level,date,logger)
);
    `
)

var (
	logTbName = "log_"
	logTbTime = "" // 201510
	logTbLock = sync.Mutex{}
)

func getLogTbName(currentTime time.Time) (string, error) {
	logTbLock.Lock()
	defer logTbLock.Unlock()
	timefmt := currentTime.Format("200601")
	// 时间相等，说明已经检查过分表名
	if logTbTime == timefmt {
		return logTbName + logTbTime, nil
	}

	tbName := logTbName + timefmt
	mdb := db.GetCache("master")
	// 创表
	if _, err := mdb.Exec(fmt.Sprintf(createLogTb, tbName)); err != nil {
		return "", errors.As(err)
	}
	// 创建成功后再刷新缓存的时间
	logTbTime = timefmt
	return tbName, nil
}

// 生成记录
func InsertLog(l *DbTable) error {
	// 获取分表名
	tbName, err := getLogTbName(l.date)
	if err != nil {
		return errors.As(err, *l)
	}
	mdb := db.GetCache("master")
	_, err = mdb.Exec(fmt.Sprintf(putLogSql, tbName),
		l.md5,
		l.platform,
		l.version,
		l.ip,
		l.date,
		l.level,
		l.logger,
		l.msg,
	)
	if err != nil {
		logTbTime = "" // 重新检查表数据是否创建了
		return errors.As(err, l)
	}
	return nil
}
