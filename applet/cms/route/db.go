package route

import (
	"strings"

	"github.com/gwaylib/database"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"

	"lserver/applet/cms/model/cms"
	"lserver/module/db"
)

var (
	mdb   = db.GetCache("master")
	cmsdb = cms.NewCmsDB()
)

func QueryDB(db *database.DB, qsql *database.Template, offset, limit int, args ...interface{}) (int, []string, [][]interface{}, error) {
	total := 0
	if err := database.QueryElem(db, &total, qsql.CountSql, args...); err != nil {
		log.Debug(errors.As(err))
		// sqlite3 table not found
		if strings.Index(err.Error(), "no such table") >= 0 {
			return 0, []string{}, [][]interface{}{}, nil
		}
		// mysql table not found
		if strings.Index(err.Error(), "Error 1146: Table") >= 0 {
			return 0, []string{}, [][]interface{}{}, nil
		}
		return 0, []string{}, [][]interface{}{}, errors.As(err, args)
	}
	if total == 0 {
		return 0, []string{}, [][]interface{}{}, nil
	}

	args = append(args, offset)
	args = append(args, limit)
	titles, datas, err := database.QueryTable(db, qsql.DataSql, args...)
	if err != nil {
		return 0, []string{}, [][]interface{}{}, errors.As(err, args)
	}
	return total, titles, datas, nil
}
