package cms

import (
	"database/sql"
	"strings"

	"lserver/module/db"

	"github.com/gwaylib/database"
	"github.com/gwaylib/errors"
)

type CmsDB interface {
	Run() error
	Close() error

	// 创建用户
	CreateUser(username, pwd, nickName string) error
	DeleteUser(username string) error
	// 更新用户状态
	UpdateUserStatus(username string, status int) error
	// 重置密码,输入原始密码
	ResetPwd(username string, pwd string) error
	// 读取用户信息
	GetUser(username string, status int) (*CmsUser, error)

	// 创建组用户
	CreateGroup(name string) (int, error)
	SetGroup(gid int, name string) error
	DeleteGroup(gid int) error

	// 菜单管理, 菜单管理用于辅助权限管理
	// 在这里仅做简单的权限设计，即只做一级uri与二级uri的菜单做权限管理
	// 例如：
	// /user
	// /user/get
	// key值格式为user, user.get
	CreateMenu(key, name string) error
	// 获取权限
	GetPriv(gid int) (CmsPriv, error)
	// 开通权限
	AddPriv(gid int, priv string) error
	// 关闭权限
	DeletePriv(gid int, priv string) error
	// 绑定模板
	BindPriv(gid int, tplname string) error
	// 开通模板权限
	AddPrivTpl(tplname string, priv string) error
	// 关闭模板权限
	DeletePrivTpl(tplname string, priv string) error
	// 删除模板
	DeleteAllPrivTpl(tplname string) error
	// 复制模板
	CopyPrivTpl(aTplName, toTplName string) error

	// 操作日志
	PutLog(username, kind, args, memo string) error
}

type cmsDB struct {
	mdb *database.DB
}

var cmsdb = NewCmsDB()

func NewCmsDB() CmsDB {
	return &cmsDB{db.GetCache("master")}
}

func (db *cmsDB) Run() error {
	return nil
}
func (db *cmsDB) Close() error {
	return db.mdb.Close()
}

const (
	createCmsUserSql = `
    INSERT INTO
        cms_user
    (
        username, passwd, nickname
    )VALUES(
        ?, ?, ?
    )
    `
)

// 创建用户
func (db *cmsDB) CreateUser(username, pwdIn, nickName string) error {
	passwd, _ := CreatePwd(pwdIn)
	_, err := db.mdb.Exec(createCmsUserSql, username, passwd, nickName)
	if err != nil {
		return errors.As(err, username, pwdIn, nickName)
	}
	return nil
}

const (
	deleteCmsUserSql = `
    DELETE FROM
        cms_user
    WHERE
        username = ?
    `
)

// 更新用户状态
func (db *cmsDB) DeleteUser(username string) error {
	_, err := db.mdb.Exec(deleteCmsUserSql, username)
	if err != nil {
		return errors.As(err, username)
	}
	return nil
}

const (
	updateCmsUserStatusSql = `
    UPDATE
        cms_user
    SET
        status = ?
    WHERE
        username = ?
    `
)

// 更新用户状态
func (db *cmsDB) UpdateUserStatus(username string, status int) error {
	_, err := db.mdb.Exec(updateCmsUserStatusSql, status, username)
	if err != nil {
		return errors.As(err, username, status)
	}
	return nil
}

const (
	updateCmsUserPwdSql = `
    UPDATE
        cms_user
    SET
        passwd = ?
    WHERE
        username = ?
    `
)

// 重置密码,输入原始密码
func (db *cmsDB) ResetPwd(username string, pwdIn string) error {
	passwd, _ := CreatePwd(pwdIn)
	_, err := db.mdb.Exec(updateCmsUserPwdSql, passwd, username)
	if err != nil {
		return errors.As(err, username, pwdIn)
	}
	return nil
}

const (
	getCmsUserSql = `
    SELECT
        passwd, nickname, gid,
        status
    FROM
        cms_user
    WHERE
        username = ?
        AND status = ?
    `
)

// 读取用户信息
func (db *cmsDB) GetUser(username string, status int) (*CmsUser, error) {
	cmsUser := &CmsUser{UserName: username}
	if err := db.mdb.QueryRow(getCmsUserSql, username, status).Scan(
		&cmsUser.Passwd,
		&cmsUser.NickName,
		&cmsUser.GroupId,
		&cmsUser.Status,
	); err != nil {
		if sql.ErrNoRows == err {
			return nil, errors.ErrNoData.As(username)
		}
		return nil, errors.As(err, username)
	}
	return cmsUser, nil
}

func (db *cmsDB) CreateGroup(name string) (int, error) {
	return 0, errors.New("TODO")
}
func (db *cmsDB) SetGroup(gid int, name string) error {
	return errors.New("TODO")
}
func (db *cmsDB) DeleteGroup(gid int) error {
	return errors.New("TODO")
}

const (
	createCmsMenuSql = `
    INSERT INTO
        cms_menu
    (
        id, name
    )VALUES(
        ?, ?
    )
    `
)

// 菜单管理, 菜单管理用于辅助权限管理
// 在这里仅做简单的权限设计，即只做一级uri与二级uri的菜单做权限管理
// 例如：
// /user
// /user/get
// key值格式为user, user.get
func (db *cmsDB) CreateMenu(key, name string) error {
	if _, err := db.mdb.Exec(createCmsMenuSql, key, name); err != nil {
		return errors.As(err, key, name)
	}
	return nil
}

const (
	getCmsMenuSql = `
    SELECT
        id, name
    FROM
        cms_menu
    LIMIT 100
    `
)

func (db *cmsDB) GetMenu() (map[string]string, error) {
	rows, err := db.mdb.Query(getCmsMenuSql)
	if err != nil {
		return nil, errors.As(err)
	}
	defer rows.Close()
	result := map[string]string{}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.As(err, id, name)
		}
		result[id] = name
	}
	return result, nil
}

const (
	getCmsPrivSql = `
    SELECT
        menu_id
    FROM
        cms_group_priv
    WHERE
        gid = ?
    ORDER BY menu_id
    `
)

// 获取权限
func (db *cmsDB) GetPriv(gid int) (CmsPriv, error) {
	rows, err := db.mdb.Query(getCmsPrivSql, gid)
	if err != nil {
		return nil, errors.As(err, gid)
	}
	defer rows.Close()
	priv := CmsPriv{}
	for rows.Next() {
		var menu_id string
		if err := rows.Scan(&menu_id); err != nil {
			return nil, errors.As(err, gid)
		}
		data := strings.Split(menu_id, ".")
		priv.Append(data)
	}
	return priv, nil
}

const (
	delCmsPrivSql = `
    DELETE FROM 
        cms_group_priv
    WHERE
        gid = ?
        AND menu_id = ?
    `
	addCmsPrivSql = `
    INSERT INTO
        cms_group_priv
    (
        gid, menu_id
    )VALUE(
        ?, ?
    )
    `
	delAllCmsPrivSql = `
    DELETE FROM 
        cms_group_priv
    WHERE
        gid = ?
    `
	bindCmsPrivSql = `
    INSERT INTO cms_group_priv(
        SELECT ?,menu_id
        FROM cms_priv_tpl
        WHERE tplname=?
    )
    `
)

// 添加权限
func (db *cmsDB) AddPriv(gid int, priv string) error {
	if _, err := db.mdb.Exec(addCmsPrivSql, gid, priv); err != nil {
		return errors.As(err, gid, priv)
	}
	return nil
}

// 关闭权限
func (db *cmsDB) DeletePriv(gid int, priv string) error {
	if _, err := db.mdb.Exec(delCmsPrivSql, gid, priv); err != nil {
		return errors.As(err, gid, priv)
	}
	return nil
}

// 绑定模板
func (db *cmsDB) BindPriv(gid int, tplname string) error {
	tx, err := db.mdb.Begin()
	if err != nil {
		return errors.As(err, gid, tplname)
	}
	if _, err := tx.Exec(delAllCmsPrivSql, gid); err != nil {
		database.Rollback(tx)
		return errors.As(err, gid, tplname)
	}
	if _, err := tx.Exec(bindCmsPrivSql, gid, tplname); err != nil {
		database.Rollback(tx)
		return errors.As(err, gid, tplname)
	}
	if err := tx.Commit(); err != nil {
		database.Rollback(tx)
		return errors.As(err, gid, tplname)
	}
	return nil
}

const (
	delCmsPrivTplSql = `
    DELETE FROM 
        cms_priv_tpl
    WHERE
        tplname = ?
        AND menu_id = ?
    `
	addCmsPrivTplSql = `
    INSERT INTO
        cms_priv_tpl
    (
        tplname, menu_id
    )VALUE(
        ?, ?
    )
    `
	delAllCmsPrivTplSql = `
    DELETE FROM 
        cms_priv_tpl
    WHERE
        tplname = ?
    `
	bindCmsPrivTplSql = `
    INSERT INTO cms_priv_tpl(
        SELECT ?,menu_id
        FROM cms_priv_tpl
        WHERE tplname=?
    )
    `
)

// 添加模板权限
func (db *cmsDB) AddPrivTpl(tplname string, priv string) error {
	if _, err := db.mdb.Exec(addCmsPrivTplSql, tplname, priv); err != nil {
		return errors.As(err, tplname, priv)
	}
	return nil
}

// 关闭模板权限
func (db *cmsDB) DeletePrivTpl(tplname string, priv string) error {
	if _, err := db.mdb.Exec(delCmsPrivTplSql, tplname, priv); err != nil {
		return errors.As(err, tplname, priv)
	}
	return nil
}

// 删除模板
func (db *cmsDB) DeleteAllPrivTpl(tplname string) error {
	if _, err := db.mdb.Exec(delAllCmsPrivTplSql, tplname); err != nil {
		return errors.As(err, tplname)
	}
	return nil
}

// 复制模板
func (db *cmsDB) CopyPrivTpl(aTplName, toTplName string) error {
	tx, err := db.mdb.Begin()
	if err != nil {
		return errors.As(err, aTplName, toTplName)
	}
	if _, err := tx.Exec(delAllCmsPrivTplSql, toTplName); err != nil {
		database.Rollback(tx)
		return errors.As(err, aTplName, toTplName)
	}
	if _, err := tx.Exec(bindCmsPrivTplSql, toTplName, aTplName); err != nil {
		database.Rollback(tx)
		return errors.As(err, aTplName, toTplName)
	}
	if err := tx.Commit(); err != nil {
		database.Rollback(tx)
		return errors.As(err, aTplName, toTplName)
	}
	return nil
}

const (
	putLogSql = `
    INSERT INTO
        cms_log
    (
        username, kind, args, memo
    )VALUES(
        ?, ?, ?, ?
    );;
    `
)

func (db *cmsDB) PutLog(username, kind, args, memo string) error {
	if _, err := db.mdb.Exec(putLogSql, username, kind, args, memo); err != nil {
		return errors.As(err, username, kind, args, memo)
	}
	return nil
}
