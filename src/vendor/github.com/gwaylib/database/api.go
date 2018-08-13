/*
此包以工厂的模式提供数据库连接，以便优化数据库连接数
*/
package database

import (
	"database/sql"
	"io"
	"runtime/debug"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
)

const (
	DRV_NAME_MYSQL     = "msyql"
	DRV_NAME_ORACLE    = "oracle" // or "oci8"
	DRV_NAME_POSTGRES  = "postgres"
	DRV_NAME_SQLITE3   = "sqlite3"
	DRV_NAME_SQLSERVER = "sqlserver" // or "mssql"
)

var (
	DEFAULT_DRV_NAME = DRV_NAME_MYSQL
)

// 使用一个已有的标准数据库实例构建出实例
func NewDB(drvName string, db *sql.DB) *DB {
	return newDB(drvName, db)
}

// 返回一个全新的实例
func Open(drvName, dsn string) (*DB, error) {
	db, err := sql.Open(drvName, dsn)
	if err != nil {
		return nil, errors.As(err, drvName, dsn)
	}
	return newDB(drvName, db), nil
}

// 注册一个池实例
func RegCache(iniFileName, sectionName string, db *DB) {
	regCache(iniFileName, sectionName, db)
}

// 获取数据库池中的实例
// 如果不存在，会使用配置文件进行读取
func GetCache(iniFileName, sectionName string) *DB {
	db, err := getCache(iniFileName, sectionName)
	if err != nil {
		panic(err)
	}
	return db
}

// 检查数据库是否存在并返回数据连接实例
func HasCache(etcFileName, sectionName string) (*DB, error) {
	return getCache(etcFileName, sectionName)
}

// 当使用了Cache，在程序退出时可调用database.CloseCache进行正常关闭数据库连接
func CloseCache() {
	closeCache()
}

// 提供懒处理的关闭方法，调用者不需要处理错误
func Close(closer io.Closer) {
	if closer == nil {
		return
	}
	if err := closer.Close(); err != nil {
		println(errors.As(err).Error())
		debug.PrintStack()
	}
}

// 提供懒处理的回滚方法，调用者不需要处理错误
func Rollback(tx *sql.Tx) {
	err := tx.Rollback()

	// roll back error is a serious error
	if err != nil {
		log.Error(errors.As(err))
	}
}

// 实现db.Exec接口
func Exec(db Execer, querySql string, args ...interface{}) (sql.Result, error) {
	return db.Exec(querySql, args...)
}

// 事务执行多个脚本
func ExecMultiTx(tx *sql.Tx, mTx []*MultiTx) error {
	return execMultiTx(tx, mTx)
}

// 通过反射添加一条数据，需要结构体至少标注字段名 `db:"name"`, 标签详情请参考github.com/jmoiron/sqlx
// 关于drvNames的设计说明
// 因支持一个可变参数, 或未填，将使用默认值:DEFAULT_DRV_NAME
func InsertStruct(exec Execer, tbName string, obj interface{}, drvNames ...string) (sql.Result, error) {
	return insertStruct(exec, tbName, obj, drvNames...)
}

// 实现db.Query查询
func Query(db Queryer, querySql string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(querySql, args...)
}

// 实现db.QueryRow查询
func QueryRow(db Queryer, querySql string, args ...interface{}) *sql.Row {
	return db.QueryRow(querySql, args...)
}

// 通过反射扫描结果至结构体
// 如果没有数据，errors.ErrNoData
func ScanStruct(rows Rows, obj interface{}) error {
	return scanStruct(rows, obj)
}

// 通过反射扫描结果至结构体数组
// 如果没有数据，返回成功，不改变原数组的值
// 代码设计请参阅github.com/jmoiron/sqlx
func ScanStructs(rows Rows, obj interface{}) error {
	return scanStructs(rows, obj)
}

// 通过反射查询结果到结构体
// 如果没有数据，返回errors.ErrNoData
func QueryStruct(db Queryer, obj interface{}, querySql string, args ...interface{}) error {
	return queryStruct(db, obj, querySql, args...)
}

// 通过反射查询多个结果到结构体数组
// 如果没有数据，返回成功，不改变原数组的值
func QueryStructs(db Queryer, obj interface{}, querySql string, args ...interface{}) error {
	return queryStructs(db, obj, querySql, args...)
}

// 查询一个支持Scan的数据类型
func QueryElem(db Queryer, result interface{}, querySql string, args ...interface{}) error {
	return queryElem(db, result, querySql, args...)
}

// 通过反射查询多个结果到数据类型数组
func QueryElems(db Queryer, result interface{}, querySql string, args ...interface{}) error {
	return queryElems(db, result, querySql, args...)
}

// 执行一个通用的查询
// 因需要查标题，相对标准sql会慢一些，适用于偷懒查询的方式
// 即使发生错误返回至少是零长度的值
func QueryTable(db Queryer, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	return queryTable(db, querySql, args...)
}

// 查询一条数据，并发map结构返回，以便页面可以直接调用
// 因需要查标题，相对标准sql会慢一些，适用于偷懒查询的方式
// 即使发生错误返回至少是零长度的值
func QueryMap(db Queryer, querySql string, args ...interface{}) ([]map[string]interface{}, error) {
	return queryMap(db, querySql, args...)
}
