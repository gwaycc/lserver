package route

import (
	"fmt"
	"strconv"

	"lserver/module/db/alarm"

	"github.com/gwaylib/database"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/eweb"
	"github.com/gwaylib/log"
	"github.com/labstack/echo"
)

const (
	LogInfoViewPath = "/app/log/info"
	LogInfoApiPath  = "/app/log/info"

	LogAlertorViewPath   = "/app/log/alertor"
	LogAlertorApiPath    = "/app/log/alertor"
	LogAlertorAddApiPath = "/app/log/alertor/add"
	LogAlertorSetApiPath = "/app/log/alertor/set"
	LogAlertorDelApiPath = "/app/log/alertor/del"

	LogMailViewPath   = "/app/log/mail"
	LogMailApiPath    = "/app/log/mail"
	LogMailSetApiPath = "/app/log/mail/set"
)

func init() {
	r := eweb.Default()
	r.GET(LogInfoViewPath, LogInfoView)
	r.POST(LogInfoApiPath, LogInfoApi)

	r.GET(LogAlertorViewPath, LogAlertorView)
	r.POST(LogAlertorApiPath, LogAlertorApi)
	r.POST(LogAlertorAddApiPath, LogAlertorAddApi)
	r.POST(LogAlertorSetApiPath, LogAlertorSetApi)
	r.POST(LogAlertorDelApiPath, LogAlertorDelApi)

	r.GET(LogMailViewPath, LogMailView)
	r.POST(LogMailApiPath, LogMailApi)
	r.POST(LogMailSetApiPath, LogMailSetApi)
}

func LogInfoView(c echo.Context) error {
	return Index(c)
}

var (
	logQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM
    %s
WHERE
    date BETWEEN ? AND ?
    AND %s
`,

		DataSql: `
SELECT
    LEFT(md5, 7) AS md5, platform, version, ip, date, level, logger, msg
FROM
    %s
WHERE
    date BETWEEN ? AND ?
    AND %s
ORDER BY date desc
LIMIT ?, ?
    `,
	}
)

func LogInfoApi(c echo.Context) error {
	currPageStr := FormValue(c, "pageId")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		log.Debug(errors.As(err, currPageStr))
		return c.String(403, "页码有误")
	}

	beginTimeStr := FormValue(c, "beginTime")
	endTimeStr := FormValue(c, "endTime")
	beginTime, err := ParseTime(beginTimeStr)
	if err != nil {
		log.Debug(errors.As(err, beginTimeStr))
		return c.String(403, "开始时间有误")
	}
	endTime, err := ParseTime(endTimeStr)
	if err != nil {
		log.Debug(errors.As(err, endTimeStr))
		return c.String(403, "结束时间有误")
	}
	tmpTime := endTime.Add(-1)
	if beginTime.Year() != tmpTime.Year() || beginTime.Month() != tmpTime.Month() {
		log.Debug(errors.As(err, endTimeStr))
		return c.String(403, "不能跨月查询")
	}
	month := beginTime.Format("200601")
	tableName := "log_" + month

	platform := FormValue(c, "platform")
	levelStr := FormValue(c, "level")
	level := 0
	logger := FormValue(c, "logger")
	msg := FormValue(c, "msg")

	var qSql *database.Template
	args := []interface{}{beginTime, endTime}
	md5Str := FormValue(c, "md5")
	if len(md5Str) == 0 {
		if len(platform) == 0 {
			return c.String(403, "请输入平台名称")
		}
		level, err = strconv.Atoi(levelStr)
		if err != nil {
			return c.String(403, "Level有误")
		}
		condition := []byte("platform = ? AND level = ?")
		args = append(args, platform)
		args = append(args, level)
		if len(logger) > 0 {
			condition = append(condition, []byte(" AND logger = ?")...)
			args = append(args, logger)
		}
		if len(msg) > 0 {
			condition = append(condition, []byte(" AND msg like ?")...)
			args = append(args, "%"+msg+"%")
		}

		qSql = logQsql.Sprintf(tableName, condition)
	} else {
		qSql = logQsql.Sprintf(tableName, "md5 like ?")
		args = append(args, md5Str+"%")
	}

	total, titles, result, err := QueryDB(
		mdb,
		qSql,
		currPage*10, 10,
		args...)
	if err != nil {
		if !errors.ErrNoData.Equal(err) {
			log.Debug(errors.As(err))
			return c.String(500, "系统错误")
		}
		// 空数据
	}
	return c.JSON(200, eweb.H{
		"total": fmt.Sprint(total),
		"names": titles,
		"data":  result,
	})
}

func LogAlertorView(c echo.Context) error {
	return Index(c)
}

func LogAlertorApi(c echo.Context) error {
	alarmCfg := alarm.NewAlarm()
	defer alarmCfg.Close()
	cfg, err := alarmCfg.LoadCfg()
	if err != nil {
		log.Info(errors.As(err))
	}

	titles := []string{"昵称", "手机", "邮件"}
	result := [][]string{}
	for _, r := range cfg.Receivers {
		result = append(result, []string{r.NickName, r.Mobile, r.Email})
	}

	return c.JSON(200, eweb.H{
		"total": fmt.Sprint(len(result)),
		"names": titles,
		"data":  result,
	})
}

func LogAlertorAddApi(c echo.Context) error {
	nickName := FormValue(c, "nickName")
	mobile := FormValue(c, "mobile")
	email := FormValue(c, "email")
	if len(nickName) == 0 {
		return c.String(403, "需要填写昵称")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(FormValue(c, "authPwd")) {
		return c.String(403, "操作密码错误")
	}
	alarmCfg := alarm.NewAlarm()
	defer alarmCfg.Close()
	cfg, err := alarmCfg.LoadCfg()
	if err != nil {
		log.Info(errors.As(err))
	}

	rs := cfg.Receivers

	_, ok := alarm.SearchReceiver(rs, nickName)
	if ok {
		return c.String(403, "昵称已存在")
	}
	cfg.Receivers = append(rs, &alarm.Receiver{NickName: nickName, Mobile: mobile, Email: email})

	if err := alarmCfg.SaveCfg(); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return nil
}

func LogAlertorSetApi(c echo.Context) error {
	nickName := FormValue(c, "nickName")
	mobile := FormValue(c, "mobile")
	email := FormValue(c, "email")
	if len(nickName) == 0 {
		return c.String(403, "需要填写昵称")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(FormValue(c, "authPwd")) {
		return c.String(403, "操作密码错误")
	}
	alarmCfg := alarm.NewAlarm()
	defer alarmCfg.Close()
	cfg, err := alarmCfg.LoadCfg()
	if err != nil {
		log.Info(errors.As(err))
	}

	rs := cfg.Receivers

	i, ok := alarm.SearchReceiver(rs, nickName)
	if !ok {
		return c.String(403, "昵称不存在")
	}
	rs[i].Mobile = mobile
	rs[i].Email = email

	if err := alarmCfg.SaveCfg(); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return nil
}

func LogAlertorDelApi(c echo.Context) error {
	nickName := FormValue(c, "nickName")
	if len(nickName) == 0 {
		return c.String(403, "需要填写昵称")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(FormValue(c, "authPwd")) {
		return c.String(403, "操作密码错误")
	}
	alarmCfg := alarm.NewAlarm()
	defer alarmCfg.Close()
	cfg, err := alarmCfg.LoadCfg()
	if err != nil {
		log.Info(errors.As(err))
	}
	rs := cfg.Receivers

	i, ok := alarm.SearchReceiver(rs, nickName)
	if !ok {
		return c.String(403, "昵称不存在")
	}
	cfg.Receivers = alarm.RemoveReceiver(rs, i)

	if err := alarmCfg.SaveCfg(); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return nil

}
func LogMailView(c echo.Context) error {
	return Index(c)
}

func LogMailApi(c echo.Context) error {
	alarmCfg := alarm.NewAlarm()
	defer alarmCfg.Close()
	cfg, err := alarmCfg.LoadCfg()
	if err != nil {
		log.Info(errors.As(err))
	}
	mailCfg := cfg.MailServer

	titles := []string{"STMP服务器", "STMP端口", "用户名", "密码"}
	result := [][]string{
		{mailCfg.SmtpHost, fmt.Sprint(mailCfg.SmtpPort), mailCfg.AuthName, "******"},
	}

	return c.JSON(200, eweb.H{
		"total": fmt.Sprint(len(result)),
		"names": titles,
		"data":  result,
	})

}

func LogMailSetApi(c echo.Context) error {
	alarmCfg := alarm.NewAlarm()
	defer alarmCfg.Close()
	cfg, err := alarmCfg.LoadCfg()
	if err != nil {
		log.Info(errors.As(err))
	}
	mailCfg := cfg.MailServer

	smtpHost := FormValue(c, "smtpHost")
	smtpPort := FormValue(c, "smtpPort")
	if len(smtpPort) == 0 {
		smtpPort = "25"
	}
	mAuthName := FormValue(c, "mAuthName")
	mAuthPwd := FormValue(c, "mAuthPwd")
	uc := GetUserCache(c)
	if !uc.ReAuth(FormValue(c, "authPwd")) {
		return c.String(403, "操作密码错误")
	}
	mailCfg.SmtpHost = smtpHost
	mailCfg.SmtpPort, err = strconv.Atoi(smtpPort)
	if err != nil {
		return c.String(403, "端口号错误")
	}
	mailCfg.AuthName = mAuthName
	mailCfg.AuthPwd = mAuthPwd
	cfg.MailServer = mailCfg

	if err := alarmCfg.Apply(cfg); err != nil {
		return c.String(403, err.Error())
	}

	if err := alarmCfg.SaveCfg(); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return nil
}
