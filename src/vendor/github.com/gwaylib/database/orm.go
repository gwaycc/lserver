package database

import (
	"github.com/jmoiron/sqlx"
)

// sqlx ORM框架
// https://github.com/jmoiron/sqlx
func (db *DB) Sqlx() *sqlx.DB {
	return db.xdb
}
