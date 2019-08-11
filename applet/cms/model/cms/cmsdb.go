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

	// 菜单管理, 菜单管理用于辅助权限管理
	// 在这里仅做简单的权限设计，即只做一级uri与二级uri的菜单做权限管理
	// 例如：
	// /user
	// /user/get
	// key值格式为user, user.get
	CreateMenu(key, name string) error
	// 获取权限
	GetPriv(username string) (CmsPriv, error)
	// 开通权限
	AddPriv(username string, priv string) error
	// 关闭权限
	DeletePriv(username string, priv string) error
	// 绑定模板
	BindPriv(username, tplname string) error
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
        passwd, nickname, group_id,
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
        cms_user_priv
    WHERE
        username = ?
    ORDER BY menu_id
    `
)

// 获取权限
func (db *cmsDB) GetPriv(username string) (CmsPriv, error) {
	rows, err := db.mdb.Query(getCmsPrivSql, username)
	if err != nil {
		return nil, errors.As(err, username)
	}
	defer rows.Close()
	priv := CmsPriv{}
	for rows.Next() {
		var menu_id string
		if err := rows.Scan(&menu_id); err != nil {
			return nil, errors.As(err, username)
		}
		data := strings.Split(menu_id, ".")
		priv.Append(data)
	}
	return priv, nil
}

const (
	delCmsPrivSql = `
    DELETE FROM 
        cms_user_priv
    WHERE
        username = ?
        AND menu_id = ?
    `
	addCmsPrivSql = `
    INSERT INTO
        cms_user_priv
    (
        username, menu_id
    )VALUE(
        ?, ?
    )
    `
	delAllCmsPrivSql = `
    DELETE FROM 
        cms_user_priv
    WHERE
        username = ?
    `
	bindCmsPrivSql = `
    INSERT INTO cms_user_priv(
        SELECT ?,menu_id
        FROM cms_user_priv_tpl
        WHERE tplname=?
    )
    `
)

// 添加权限
func (db *cmsDB) AddPriv(username string, priv string) error {
	if _, err := db.mdb.Exec(addCmsPrivSql, username, priv); err != nil {
		return errors.As(err, username, priv)
	}
	return nil
}

// 关闭权限
func (db *cmsDB) DeletePriv(username string, priv string) error {
	if _, err := db.mdb.Exec(delCmsPrivSql, username, priv); err != nil {
		return errors.As(err, username, priv)
	}
	return nil
}

// 绑定模板
func (db *cmsDB) BindPriv(username, tplname string) error {
	tx, err := db.mdb.Begin()
	if err != nil {
		return errors.As(err, username, tplname)
	}
	if _, err := tx.Exec(delAllCmsPrivSql, username); err != nil {
		database.Rollback(tx)
		return errors.As(err, username, tplname)
	}
	if _, err := tx.Exec(bindCmsPrivSql, username, tplname); err != nil {
		database.Rollback(tx)
		return errors.As(err, username, tplname)
	}
	if err := tx.Commit(); err != nil {
		database.Rollback(tx)
		return errors.As(err, username, tplname)
	}
	return nil
}

const (
	delCmsPrivTplSql = `
    DELETE FROM 
        cms_user_priv_tpl
    WHERE
        tplname = ?
        AND menu_id = ?
    `
	addCmsPrivTplSql = `
    INSERT INTO
        cms_user_priv_tpl
    (
        tplname, menu_id
    )VALUE(
        ?, ?
    )
    `
	delAllCmsPrivTplSql = `
    DELETE FROM 
        cms_user_priv_tpl
    WHERE
        tplname = ?
    `
	bindCmsPrivTplSql = `
    INSERT INTO cms_user_priv_tpl(
        SELECT ?,menu_id
        FROM cms_user_priv_tpl
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
