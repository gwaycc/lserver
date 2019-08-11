package route

import (
	"fmt"
	"strconv"

	"lserver/applet/cms/model/cms"

	"github.com/gwaylib/database"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/eweb"
	"github.com/gwaylib/log"
	"github.com/labstack/echo"
)

const (
	CmsInfoViewPath  = "/app/cms/info"
	CmsInfoApiPath   = "/app/cms/info"
	CmsCreateApiPath = "/app/cms/create"
	CmsPwdApiPath    = "/app/cms/pwd"
	CmsDeleteApiPath = "/app/cms/delete"

	CmsLogViewPath = "/app/cms/log"
	CmsLogApiPath  = "/app/cms/log"

	CmsPrivViewPath    = "/app/cms/priv"
	CmsPrivApiPath     = "/app/cms/priv"
	CmsPrivBindApiPath = "/app/cms/priv/bind"
	CmsPrivOnApiPath   = "/app/cms/priv/on"
	CmsPrivOffApiPath  = "/app/cms/priv/off"

	CmsPrivTplViewPath      = "/app/cms/privtpl"
	CmsPrivTplApiPath       = "/app/cms/privtpl"
	CmsPrivListTplApiPath   = "/app/cms/priv/tpl/list"
	CmsPrivOnTplApiPath     = "/app/cms/priv/tpl/on"
	CmsPrivOffTplApiPath    = "/app/cms/priv/tpl/off"
	CmsPrivNewTplApiPath    = "/app/cms/priv/tpl/new"
	CmsPrivDeleteTplApiPath = "/app/cms/priv/tpl/delete"
)

func init() {
	r := eweb.Default()
	r.GET(CmsInfoViewPath, CmsInfoView)
	r.POST(CmsInfoApiPath, CmsInfoApi)
	r.POST(CmsCreateApiPath, CmsCreateApi)
	r.POST(CmsPwdApiPath, CmsPwdApi)
	r.POST(CmsDeleteApiPath, CmsDeleteApi)

	r.GET(CmsLogViewPath, CmsLogView)
	r.POST(CmsLogApiPath, CmsLogApi)

	r.GET(CmsPrivViewPath, CmsPrivView)
	r.POST(CmsPrivApiPath, CmsPrivApi)
	r.POST(CmsPrivBindApiPath, CmsPrivBindApi)
	r.POST(CmsPrivOnApiPath, CmsPrivOnApi)
	r.POST(CmsPrivOffApiPath, CmsPrivOffApi)

	r.GET(CmsPrivTplViewPath, CmsPrivTplView)
	r.POST(CmsPrivTplApiPath, CmsPrivTplApi)
	r.POST(CmsPrivListTplApiPath, CmsPrivListTplApi)
	r.POST(CmsPrivOnTplApiPath, CmsPrivOnTplApi)
	r.POST(CmsPrivOffTplApiPath, CmsPrivOffTplApi)
	r.POST(CmsPrivNewTplApiPath, CmsPrivNewTplApi)
	r.POST(CmsPrivDeleteTplApiPath, CmsPrivDeleteTplApi)
}

func CmsInfoView(c echo.Context) error {
	return Index(c)
}

var (
	cmsInfoQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM
    cms_user
WHERE
    username like ?
    `,

		DataSql: `
SELECT 
    username "帐号",
    nickname "昵称",
    CASE status WHEN 1 THEN '可用' WHEN 2 THEN '禁用' ELSE status END "状态",
    created_at "创建时间" 
FROM
    cms_user
WHERE
    status = 1
    AND username like ?
ORDER BY username
LIMIT ?, ?
    `,
	}
)

func CmsInfoApi(c echo.Context) error {
	currPageStr := FormValue(c, "pageId")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		log.Debug(errors.As(err, currPageStr))
		return c.String(403, "页码有误")
	}

	userName := FormValue(c, "userName")
	total, titles, result, err := QueryDB(
		mdb,
		cmsInfoQsql,
		currPage*10, 10,
		"%"+userName+"%")
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

func CmsCreateApi(c echo.Context) error {
	userName := FormValue(c, "userName")
	userPwd := FormValue(c, "userPwd")
	nickName := FormValue(c, "nickName")
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if len(userName) == 0 {
		log.Debug("no userName")
		return c.String(403, "请输入用户帐号")
	}
	if len(userPwd) == 0 {
		log.Debug("no userPwd")
		return c.String(403, "请输入用户密码")
	}
	if len(nickName) == 0 {
		log.Debug("no nickName")
		return c.String(403, "请输入用户昵称")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	cmsdb := cms.NewCmsDB()
	if _, err := cmsdb.GetUser(userName, 1); err == nil {
		return c.String(403, "用户已存在")
	} else if !errors.ErrNoData.Equal(err) {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	if err := cmsdb.CreateUser(userName, userPwd, nickName); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}

	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "创建后台用户", userName, memo); err != nil {
		log.Warn(errors.As(err))
	}

	return c.String(200, "操作成功")
}
func CmsPwdApi(c echo.Context) error {
	userName := FormValue(c, "userName")
	userPwd := FormValue(c, "userPwd")
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if len(userName) == 0 {
		log.Debug("no userName")
		return c.String(403, "请输入用户帐号")
	}
	if len(userPwd) == 0 {
		log.Debug("no userPwd")
		return c.String(403, "请输入用户密码")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	if err := cmsdb.ResetPwd(userName, userPwd); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "修改后台密码", userName, memo); err != nil {
		log.Warn(errors.As(err))
	}

	return c.String(200, "操作成功")
}
func CmsDeleteApi(c echo.Context) error {
	userName := FormValue(c, "userName")
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if len(userName) == 0 {
		log.Debug("no userName")
		return c.String(403, "请输入用户帐号")
	}
	if userName == "admin" {
		return c.String(403, "内置帐号不支持删除")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.DeleteUser(userName); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "删除后台用户", userName, memo); err != nil {
		log.Warn(errors.As(err))
	}

	return c.String(200, "操作成功")
}

func CmsLogView(c echo.Context) error {
	return Index(c)
}

var (
	cmsLogQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM
    cms_log
WHERE
    created_at BETWEEN ? AND ?
    AND username like ?

    `,

		DataSql: `
SELECT 
    created_at "操作时间",
    username "操作人",
    kind "操作类别",
    args "输入参数",
    memo "备注"
FROM
    cms_log
WHERE
    created_at BETWEEN ? AND ?
    AND username like ?
ORDER BY created_at DESC
LIMIT ?, ?
    `,
	}
)

func CmsLogApi(c echo.Context) error {
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
	userName := FormValue(c, "userName")
	total, titles, result, err := QueryDB(
		mdb,
		cmsLogQsql,
		currPage*10, 10,
		beginTime, endTime, "%"+userName+"%")
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

func CmsPrivView(c echo.Context) error {
	return Index(c)
}

var (
	cmsPrivQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM cms_menu tb1
LEFT JOIN cms_user_priv tb2 on tb1.id=tb2.menu_id and tb2.username=?
WHERE
    %s
	AND tb1.name like ?
    AND tb1.id <> '*.*.*.*.*'
    `,

		DataSql: `
SELECT 
    tb1.id "功能ID", tb1.name "功能名称", CASE WHEN tb2.menu_id IS NOT NULL THEN '已开通' ELSE '未开通' END '是否开通'
FROM cms_menu tb1
LEFT JOIN cms_user_priv tb2 on tb1.id=tb2.menu_id and tb2.username=?
WHERE
    %s
	AND tb1.name like ?
    AND tb1.id <> '*.*.*.*.*'
ORDER BY tb1.name
LIMIT ?, ?
    `,
	}
)

func CmsPrivApi(c echo.Context) error {
	currPageStr := FormValue(c, "pageId")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		log.Debug(errors.As(err, currPageStr))
		return c.String(403, "页码有误")
	}

	userName := FormValue(c, "userName")
	if len(userName) == 0 {
		return c.String(403, "请输入用户帐号")
	}
	menuName := FormValue(c, "menuName")
	// 校验帐号是否存在
	cmsdb := cms.NewCmsDB()
	if _, err := cmsdb.GetUser(userName, 1); err != nil {
		if errors.ErrNoData.Equal(err) {
			log.Debug(errors.As(err, userName))
			return c.String(403, "用户名不存在")
		}
		log.Warn(errors.As(err, userName))
		return c.String(500, "系统错误")
	}

	status := FormValue(c, "status")
	where := "1=1"
	switch status {
	case "1":
		where = "tb2.username IS NOT NULL"
	case "2":
		where = "tb2.username IS NULL"
	}
	total, titles, result, err := QueryDB(
		mdb,
		cmsPrivQsql.Sprintf(where),
		currPage*10, 10,
		userName, "%"+menuName+"%")
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

func CmsPrivBindApi(c echo.Context) error {
	userName := FormValue(c, "userName")
	if len(userName) == 0 {
		return c.String(403, "请输入用户帐号")
	}
	if userName == "admin" {
		return c.String(403, "内置帐号不支持设定")
	}
	tplName := FormValue(c, "tplName")
	if len(tplName) == 0 {
		return c.String(403, "请输入功能模板")
	}
	cmsdb := cms.NewCmsDB()
	// 校验用户是否存在
	if _, err := cmsdb.GetUser(userName, 1); err != nil {
		if errors.ErrNoData.Equal(err) {
			return c.String(403, "帐号不存在")
		}
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	if err := cmsdb.BindPriv(userName, tplName); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivOnApi(c echo.Context) error {
	userName := FormValue(c, "userName")
	if len(userName) == 0 {
		return c.String(403, "请输入用户帐号")
	}
	menuId := FormValue(c, "menuId")
	if len(menuId) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.AddPriv(userName, menuId); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivOffApi(c echo.Context) error {
	userName := FormValue(c, "userName")
	if len(userName) == 0 {
		return c.String(403, "请输入用户帐号")
	}
	menuId := FormValue(c, "menuId")
	if len(menuId) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.DeletePriv(userName, menuId); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivTplView(c echo.Context) error {
	return Index(c)
}

var (
	cmsPrivTplQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM cms_menu tb1
LEFT JOIN cms_user_priv_tpl tb2 on tb1.id=tb2.menu_id and tplname=?
WHERE
    %s
	AND tb1.name like ?
    AND tb1.id <> '*.*.*.*.*'
    `,

		DataSql: `
SELECT 
    tb1.id "功能ID", tb1.name "功能名称", CASE WHEN tb2.menu_id IS NOT NULL THEN '已开通' ELSE '未开通' END '是否开通'
FROM 
    cms_menu tb1
    LEFT JOIN cms_user_priv_tpl tb2 on tb1.id=tb2.menu_id and tb2.tplname=?
WHERE
    %s
	AND tb1.name like ?
    AND tb1.id <> '*.*.*.*.*'
ORDER BY tb1.name
LIMIT ?, ?
    `,
	}
)

func CmsPrivTplApi(c echo.Context) error {
	currPageStr := FormValue(c, "pageId")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		log.Debug(errors.As(err, currPageStr))
		return c.String(403, "页码有误")
	}

	tplName := FormValue(c, "tplName")
	if len(tplName) == 0 {
		return c.String(403, "请输入模板名称")
	}
	menuName := FormValue(c, "menuName")
	status := FormValue(c, "status")
	where := "1=1"
	switch status {
	case "1":
		where = "tb2.menu_id IS NOT NULL"
	case "2":
		where = "tb2.menu_id IS NULL"
	}

	total, titles, result, err := QueryDB(
		mdb,
		cmsPrivTplQsql.Sprintf(where),
		currPage*10, 10,
		tplName, "%"+menuName+"%")
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

var (
	cmsPrivListTplQsql = &database.Template{
		CountSql: `
SELECT 
    count(1)
FROM
    cms_user_priv_tpl
GROUP BY tplname
ORDER BY tplname
		`,
		DataSql: `
SELECT 
    tplname
FROM
    cms_user_priv_tpl
GROUP BY tplname
ORDER BY tplname
LIMIT ?, ?
    `,
	}
)

func CmsPrivListTplApi(c echo.Context) error {
	// 仅支持10个模板，若多了，需要重新设计页面显示
	total, titles, result, err := QueryDB(
		mdb,
		cmsPrivListTplQsql,
		0, 10)
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

func CmsPrivOnTplApi(c echo.Context) error {
	tplName := FormValue(c, "tplName")
	if len(tplName) == 0 {
		return c.String(403, "请输入模板名称")
	}
	menuId := FormValue(c, "menuId")
	if len(menuId) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.AddPrivTpl(tplName, menuId); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivOffTplApi(c echo.Context) error {
	tplName := FormValue(c, "tplName")
	if len(tplName) == 0 {
		return c.String(403, "请输入模板名称")
	}
	menuId := FormValue(c, "menuId")
	if len(menuId) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.DeletePrivTpl(tplName, menuId); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivNewTplApi(c echo.Context) error {
	aTplName := FormValue(c, "aTplName")
	if len(aTplName) == 0 {
		return c.String(403, "请输入来源模板名称")
	}
	toTplName := FormValue(c, "toTplName")
	if len(toTplName) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.CopyPrivTpl(aTplName, toTplName); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivDeleteTplApi(c echo.Context) error {
	tplName := FormValue(c, "tplName")
	if len(tplName) == 0 {
		return c.String(403, "请输入模板名称")
	}

	cmsdb := cms.NewCmsDB()
	if err := cmsdb.DeleteAllPrivTpl(tplName); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}
