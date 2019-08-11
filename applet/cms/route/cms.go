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
	CmsUserInfoViewPath  = "/app/cms/user/info"
	CmsUserInfoApiPath   = "/app/cms/user/info"
	CmsUserCreateApiPath = "/app/cms/user/create"
	CmsUserPwdApiPath    = "/app/cms/user/pwd"
	CmsUserGroupApiPath  = "/app/cms/user/group"
	CmsUserDeleteApiPath = "/app/cms/user/delete"

	CmsGroupInfoViewPath  = "/app/cms/group/info"
	CmsGroupInfoApiPath   = "/app/cms/group/info"
	CmsGroupCreateApiPath = "/app/cms/group/create"
	CmsGroupSetApiPath    = "/app/cms/group/set"
	CmsGroupDeleteApiPath = "/app/cms/group/delete"

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
	r.GET(CmsUserInfoViewPath, CmsUserInfoView)
	r.POST(CmsUserInfoApiPath, CmsUserInfoApi)
	r.POST(CmsUserCreateApiPath, CmsUserCreateApi)
	r.POST(CmsUserPwdApiPath, CmsUserPwdApi)
	r.POST(CmsUserGroupApiPath, CmsUserGroupApi)
	r.POST(CmsUserDeleteApiPath, CmsUserDeleteApi)

	r.GET(CmsGroupInfoViewPath, CmsGroupInfoView)
	r.POST(CmsGroupInfoApiPath, CmsGroupInfoApi)
	r.POST(CmsGroupCreateApiPath, CmsGroupCreateApi)
	r.POST(CmsGroupSetApiPath, CmsGroupSetApi)
	r.POST(CmsGroupDeleteApiPath, CmsGroupDeleteApi)

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

func CmsUserInfoView(c echo.Context) error {
	return Index(c)
}

var (
	cmsUserInfoQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM
    cms_user
WHERE
    status = 1
    AND username like ?
    `,

		DataSql: `
SELECT 
    tb1.username "帐号",
    tb1.nickname "昵称",
    tb1.gid "组ID",
    CASE WHEN tb2.name IS NULL THEN '<未知>' ELSE tb2.name END "组名称",
    CASE tb1.status WHEN 1 THEN '可用' WHEN 2 THEN '禁用' ELSE tb1.status END "状态",
    tb1.created_at "创建时间" 
FROM
    cms_user tb1
	LEFT JOIN cms_group tb2 ON tb1.gid=tb2.id
WHERE
    tb1.status = 1
    AND tb1.username like ?
ORDER BY tb1.username
LIMIT ?, ?
    `,
	}
)

func CmsUserInfoApi(c echo.Context) error {
	currPageStr := FormValue(c, "pageId")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		log.Debug(errors.As(err, currPageStr))
		return c.String(403, "页码有误")
	}

	userName := FormValue(c, "userName")
	total, titles, result, err := QueryDB(
		mdb,
		cmsUserInfoQsql,
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

func CmsUserCreateApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID错误")
	}
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
	pwd, _ := cms.CreatePwd(userPwd)
	if err := cmsdb.CreateUser(&cms.CmsUser{
		UserName: userName,
		Passwd:   pwd,
		NickName: nickName,
		Gid:      gid,
	}); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}

	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "创建后台用户", userName, memo); err != nil {
		log.Warn(errors.As(err))
	}

	return c.String(200, "操作成功")
}

func CmsUserPwdApi(c echo.Context) error {
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
func CmsUserGroupApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID不正确")
	}
	userName := FormValue(c, "userName")
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if len(userName) == 0 {
		log.Debug("no userName")
		return c.String(403, "请输入用户帐号")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	if err := cmsdb.UpdateUserGroup(userName, gid); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "修改用户组", fmt.Sprint(userName, gid), memo); err != nil {
		log.Warn(errors.As(err))
	}

	return c.String(200, "操作成功")
}
func CmsUserDeleteApi(c echo.Context) error {
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

func CmsGroupInfoView(c echo.Context) error {
	return Index(c)
}

var (
	cmsGroupInfoQsql = &database.Template{
		CountSql: `
SELECT 
    count(1) 
FROM
    cms_group
WHERE
    name like ?
    `,

		DataSql: `
SELECT 
    tb1.id "组ID",
    tb1.name "组名称",
    tb1.created_at "创建时间" 
FROM
    cms_group tb1
WHERE
    tb1.name like ?
LIMIT ?, ?
    `,
	}
)

func CmsGroupInfoApi(c echo.Context) error {
	currPageStr := FormValue(c, "pageId")
	currPage, err := strconv.Atoi(currPageStr)
	if err != nil {
		log.Debug(errors.As(err, currPageStr))
		return c.String(403, "页码有误")
	}

	name := FormValue(c, "name")
	total, titles, result, err := QueryDB(
		mdb,
		cmsGroupInfoQsql,
		currPage*10, 10,
		"%"+name+"%")
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

func CmsGroupCreateApi(c echo.Context) error {
	name := FormValue(c, "name")
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if len(name) == 0 {
		log.Debug("no name")
		return c.String(403, "请输入组名称")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	cmsdb := cms.NewCmsDB()
	if _, err := cmsdb.CreateGroup(name); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}

	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "创建后台分组", name, memo); err != nil {
		log.Warn(errors.As(err))
	}
	return c.String(200, "操作成功")
}

func CmsGroupSetApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID不正确")
	}
	name := FormValue(c, "name")
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if len(name) == 0 {
		log.Debug("no name")
		return c.String(403, "请输入组名称")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	if err := cmsdb.SetGroup(gid, name); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "修改用户组", fmt.Sprint(gid, name), memo); err != nil {
		log.Warn(errors.As(err))
	}

	return c.String(200, "操作成功")
}
func CmsGroupDeleteApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID错误")
	}
	memo := FormValue(c, "memo")
	authPwd := FormValue(c, "authPwd")
	if gid == 0 {
		return c.String(403, "内置分组不支持删除")
	}
	uc := GetUserCache(c)
	if !uc.ReAuth(authPwd) {
		log.Debug(authPwd)
		return c.String(403, "操作密码错误")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.DeleteGroup(gid); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	// 生成日志
	if err := cmsdb.PutLog(uc.UserName, "删除后台分组", fmt.Sprint(gid), memo); err != nil {
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
LEFT JOIN cms_group_priv tb2 on tb1.id=tb2.menu_id and tb2.gid=?
WHERE
    %s
	AND tb1.name like ?
    AND tb1.id <> '*.*.*.*.*'
    `,

		DataSql: `
SELECT 
    tb1.id "功能ID", tb1.name "功能名称", CASE WHEN tb2.menu_id IS NOT NULL THEN '已开通' ELSE '未开通' END '是否开通'
FROM cms_menu tb1
LEFT JOIN cms_group_priv tb2 on tb1.id=tb2.menu_id and tb2.gid=?
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
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID不正确")
	}
	menuName := FormValue(c, "menuName")
	// 校验帐号是否存在
	cmsdb := cms.NewCmsDB()
	if _, err := cmsdb.GetGroup(gid); err != nil {
		if errors.ErrNoData.Equal(err) {
			log.Debug(errors.As(err, gid))
			return c.String(403, "组ID不存在")
		}
		log.Warn(errors.As(err, gid))
		return c.String(500, "系统错误")
	}

	status := FormValue(c, "status")
	where := "1=1"
	switch status {
	case "1":
		where = "tb2.gid IS NOT NULL"
	case "2":
		where = "tb2.gid IS NULL"
	}
	total, titles, result, err := QueryDB(
		mdb,
		cmsPrivQsql.Sprintf(where),
		currPage*10, 10,
		gid, "%"+menuName+"%")
	if err != nil {
		log.Debug(errors.As(err))
		return c.String(500, "系统错误")
		// 空数据
	}
	return c.JSON(200, eweb.H{
		"total": fmt.Sprint(total),
		"names": titles,
		"data":  result,
	})
}

func CmsPrivBindApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID不正确")
	}
	if gid == 0 {
		return c.String(403, "内置组不支持设定")
	}
	tplName := FormValue(c, "tplName")
	if len(tplName) == 0 {
		return c.String(403, "请输入功能模板")
	}
	cmsdb := cms.NewCmsDB()
	// 校验是否存在
	if _, err := cmsdb.GetGroup(gid); err != nil {
		if errors.ErrNoData.Equal(err) {
			return c.String(403, "组不存在")
		}
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	if err := cmsdb.BindPriv(gid, tplName); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivOnApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID不正确")
	}
	if gid == 0 {
		return c.String(403, "内置组不支持设定")
	}
	menuId := FormValue(c, "menuId")
	if len(menuId) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	// 校验是否存在
	if _, err := cmsdb.GetGroup(gid); err != nil {
		if errors.ErrNoData.Equal(err) {
			return c.String(403, "组不存在")
		}
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}

	if err := cmsdb.AddPriv(gid, menuId); err != nil {
		log.Warn(errors.As(err))
		return c.String(500, "系统错误")
	}
	return c.String(200, "操作成功")
}

func CmsPrivOffApi(c echo.Context) error {
	gid, err := strconv.Atoi(FormValue(c, "gid"))
	if err != nil {
		return c.String(403, "组ID不正确")
	}
	if gid == 0 {
		return c.String(403, "内置组不支持设定")
	}
	menuId := FormValue(c, "menuId")
	if len(menuId) == 0 {
		return c.String(403, "请输入菜单编号")
	}
	cmsdb := cms.NewCmsDB()
	if err := cmsdb.DeletePriv(gid, menuId); err != nil {
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
LEFT JOIN cms_priv_tpl tb2 on tb1.id=tb2.menu_id and tplname=?
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
    LEFT JOIN cms_priv_tpl tb2 on tb1.id=tb2.menu_id and tb2.tplname=?
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
    cms_priv_tpl
GROUP BY tplname
ORDER BY tplname
		`,
		DataSql: `
SELECT 
    tplname
FROM
    cms_priv_tpl
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
