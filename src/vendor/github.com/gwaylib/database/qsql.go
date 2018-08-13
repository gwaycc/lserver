package database

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gwaylib/errors"
	"github.com/jmoiron/sqlx/reflectx"
)

// 自增回调接口
type AutoIncrAble interface {
	// notify for last id
	SetLastInsertId(id int64, err error)
}

// 执行器
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// 查询器
type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// 扫描器
type Rows interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
}

type MultiTx struct {
	Query string
	Args  []interface{}
}

func NewMultiTx(query string, args ...interface{}) *MultiTx {
	return &MultiTx{query, args}
}

const (
	addObjSql = "INSERT INTO %s (%s) VALUES (%s);"
)

// 添加一条数据，需要结构体至少标注字段名 `db:"name"`, 标签详情请参考github.com/jmoiron/sqlx
// 关于drvNames的设计说明
// 因支持一个可变参数, 或未填，将使用默认值:DEFAULT_DRV_NAME
func insertStruct(exec Execer, tbName string, obj interface{}, drvNames ...string) (sql.Result, error) {
	// 自动检查数据库驱动名称
	drvName := DEFAULT_DRV_NAME
	db, ok := exec.(*DB)
	if ok {
		drvName = db.DriverName()
	} else {
		drvNamesLen := len(drvNames)
		if drvNamesLen > 0 {
			if drvNamesLen != 0 {
				panic(errors.New("'drvNames' expect only one argument").As(drvNames))
			}
			drvName = drvNames[0]
		}
	}

	names, inputs, vals, err := reflectInsertStruct(obj, drvName)
	if err != nil {
		return nil, errors.As(err)
	}
	execSql := fmt.Sprintf(addObjSql, tbName, names, inputs)
	// log.Debugf("%s%+v", execSql, vals)
	result, err := exec.Exec(execSql, vals...)
	if err != nil {
		return nil, errors.As(err, execSql)
	}
	incr, ok := obj.(AutoIncrAble) // need obj is ptr kind.
	if ok {
		incr.SetLastInsertId(result.LastInsertId())
	}
	return result, nil
}

func execMultiTx(tx *sql.Tx, mTx []*MultiTx) error {
	for _, mt := range mTx {
		if _, err := tx.Exec(mt.Query, mt.Args...); err != nil {
			return errors.As(err)
		}
	}
	return nil
}

// fieldsByName fills a values interface with fields from the passed value based
// on the traversals in int.  If ptrs is true, return addresses instead of values.
// We write this instead of using FieldsByName to save allocations and map lookups
// when iterating over many rows.  Empty traversals will get an interface pointer.
// Because of the necessity of requesting ptrs or values, it's considered a bit too
// specialized for inclusion in reflectx itself.
func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}
	return nil
}

func scanStruct(rows Rows, obj interface{}) error {
	if obj == nil {
		return errors.New("nil pointer passed to StructScan destination")
	}

	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	base := reflectx.Deref(value.Type())
	if base.Kind() != reflect.Struct {
		return errors.As(fmt.Errorf("expected struct pointer but got %s", value.Kind()))
	}

	columns, err := rows.Columns()
	if err != nil {
		return errors.As(err)
	}

	fields := refxM.TraversalsByName(base, columns)
	values := make([]interface{}, len(columns))

	direct := reflect.Indirect(value)

	vp := reflect.New(base)
	v := reflect.Indirect(vp)
	if err := fieldsByTraversal(v, fields, values, true); err != nil {
		return errors.As(err)
	}
	if !rows.Next() {
		return errors.ErrNoData
	}
	if err := rows.Scan(values...); err != nil {
		return errors.As(err)
	}
	direct.Set(v)
	return nil
}
func scanStructs(rows Rows, obj interface{}) error {
	if obj == nil {
		return errors.New("nil pointer passed to StructScan destination")
	}
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	slice := reflectx.Deref(value.Type())
	if slice.Kind() != reflect.Slice {
		return errors.As(fmt.Errorf("expected slice but got %s", value.Kind()))
	}
	base := reflectx.Deref(slice.Elem())

	columns, err := rows.Columns()
	if err != nil {
		return errors.As(err)
	}
	fields := refxM.TraversalsByName(base, columns)
	direct := reflect.Indirect(value)
	isPtr := slice.Elem().Kind() == reflect.Ptr
	values := make([]interface{}, len(columns))
	var v, vp reflect.Value
	for rows.Next() {
		vp = reflect.New(base)
		v = reflect.Indirect(vp)
		if err := fieldsByTraversal(v, fields, values, true); err != nil {
			return errors.As(err)
		}

		if err := rows.Scan(values...); err != nil {
			return errors.As(err)
		}
		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, v))
		}
	}

	return nil
}

func queryStruct(db Queryer, obj interface{}, querySql string, args ...interface{}) error {
	rows, err := db.Query(querySql, args...)
	if err != nil {
		return errors.As(err, args)
	}
	defer Close(rows)

	if err := scanStruct(rows, obj); err != nil {
		return errors.As(err, args)
	}
	return nil
}

func queryStructs(db Queryer, obj interface{}, querySql string, args ...interface{}) error {
	rows, err := db.Query(querySql, args...)
	if err != nil {
		return errors.As(err, args)
	}
	defer Close(rows)

	if err := scanStructs(rows, obj); err != nil {
		return errors.As(err, args)
	}

	return nil
}

// 查询一个支持Scan的数据类型
func queryElem(db Queryer, result interface{}, querySql string, args ...interface{}) error {
	if err := db.QueryRow(querySql, args...).Scan(result); err != nil {
		if sql.ErrNoRows != err {
			return errors.As(err, querySql, args)
		}
		return errors.ErrNoData.As(args)
	}
	return nil
}

// 查询一个支持Scan的数据类型数组
func queryElems(db Queryer, arr interface{}, querySql string, args ...interface{}) error {
	if arr == nil {
		return errors.New("nil pointer passed to StructScan destination")
	}
	value := reflect.ValueOf(arr)
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	slice := reflectx.Deref(value.Type())
	if slice.Kind() != reflect.Slice {
		return errors.As(fmt.Errorf("expected slice but got %s", value.Kind()))
	}
	base := reflectx.Deref(slice.Elem())

	rows, err := db.Query(querySql, args...)
	if err != nil {
		return errors.As(err, querySql, args)
	}
	defer Close(rows)

	isPtr := slice.Elem().Kind() == reflect.Ptr
	direct := reflect.Indirect(value)
	var vp reflect.Value
	for rows.Next() {
		vp = reflect.New(base)
		if err := rows.Scan(vp.Interface()); err != nil {
			return errors.As(err)
		}
		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}
	return nil
}

// 执行一个通用的查询
// 因需要查标题，相对标准sql会慢一些，适用于偷懒查询的方式
// 即使发生错误返回至少是零长度的值
func queryTable(db Queryer, querySql string, args ...interface{}) (titles []string, result [][]interface{}, err error) {
	titles = []string{}
	result = [][]interface{}{}
	rows, err := db.Query(querySql, args...)
	if err != nil {
		return titles, result, errors.As(err, args)
	}
	defer Close(rows)

	titles, err = rows.Columns()
	if err != nil {
		return titles, result, errors.As(err, args)
	}

	for rows.Next() {
		r := makeDBData(len(titles))
		if err := rows.Scan(r...); err != nil {
			return titles, result, errors.As(err, args)
		}
		result = append(result, r)
	}
	if len(result) == 0 {
		return titles, result, errors.ErrNoData.As(args)
	}

	return titles, result, nil
}

// 查询一条数据，并发map结构返回，以便页面可以直接调用
// 因需要查标题，相对标准sql会慢一些，适用于偷懒查询的方式
// 即使发生错误返回至少是零长度的值
func queryMap(db Queryer, querySql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(querySql, args...)
	if err != nil {
		return []map[string]interface{}{}, errors.As(err, args)
	}
	defer Close(rows)

	// 列名
	names, err := rows.Columns()
	if err != nil {
		return []map[string]interface{}{}, errors.As(err, args)
	}

	// 取一条数据
	result := []map[string]interface{}{}
	for rows.Next() {
		r := makeDBData(len(names))
		if err := rows.Scan(r...); err != nil {
			return []map[string]interface{}{}, errors.As(err, args)
		}
		mData := map[string]interface{}{}
		for i, name := range names {
			// 校验列名重复性
			_, ok := mData[name]
			if ok {
				return result, errors.New("Already exist column name").As(name)
			}
			mData[name] = r[i]
		}
		result = append(result, mData)
	}
	if len(result) == 0 {
		return []map[string]interface{}{}, errors.ErrNoData.As(err, args)
	}
	return result, nil
}
