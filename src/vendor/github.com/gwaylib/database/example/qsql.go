package main

import (
	"fmt"

	"github.com/gwaylib/database"
	"github.com/gwaylib/errors"
	_ "github.com/mattn/go-sqlite3"
)

type TestingUser struct {
	ID       int64  `db:"id,auto_increment"`
	UserName string `db:"username"`
	Passwd   string `db:"passwd"`
}

// implement LastInsertId callback
func (u *TestingUser) SetLastInsertId(id int64, err error) {
	if err != nil {
		panic(err)
	}
	u.ID = id
}

func main() {
	mdb, _ := database.Open("sqlite3", ":memory:")
	defer database.Close(mdb)

	// create table
	if _, err := database.Exec(mdb,
		`CREATE TABLE user (
		  "id" INTEGER PRIMARY KEY NOT NULL,
		  "username" VARCHAR(32) NOT NULL UNIQUE,
		  "passwd" VARCHAR(128) NOT NULL,
		  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`); err != nil {
		panic(err)
	}

	// std insert
	if _, err := database.Exec(mdb, "INSERT INTO user(username,passwd)VALUES(?,?)", "t1", "t1"); err != nil {
		panic(err)
	}

	// reflect insert
	newUser := &TestingUser{UserName: "t2", Passwd: "t2"}
	if _, err := database.InsertStruct(mdb, "user", newUser); err != nil {
		panic(err)
	}
	if newUser.ID == 0 {
		panic("expect newUser.ID > 0")
	}

	// std query
	var id int64
	var username, passwd string
	if err := database.QueryRow(mdb, "SELECT id, username, passwd FROM user WHERE username=?", "t1").Scan(&id, &username, &passwd); err != nil {
		panic(err)
	}
	if username != "t1" && passwd != "t1" {
		panic(username + "," + passwd)
	}
	if id == 0 {
		panic(id)
	}

	// reflect query
	// query struct data
	expectUser := &TestingUser{}
	if err := database.QueryStruct(mdb, expectUser, "SELECT * FROM user WHERE username=?", "t1"); err != nil {
		panic(err)
	}
	if expectUser.UserName != "t1" && expectUser.Passwd != "t1" {
		panic("data not match")
	}
	users := []*TestingUser{}
	if err := database.QueryStructs(mdb, &users, "SELECT * FROM user LIMIT 2"); err != nil {
		panic(err)
	}
	if len(users) != 2 {
		panic("expect len==2")
	}
	// query elememt data
	pwd := ""
	if err := database.QueryElem(mdb, &pwd, "SELECT passwd FROM user WHERE username=?", "t1"); err != nil {
		panic(err)
	}
	if pwd != "t1" {
		panic(pwd)
	}
	pwds := []string{}
	if err := database.QueryElems(mdb, &pwds, "SELECT passwd FROM user LIMIT 2"); err != nil {
		panic(err)
	}
	if len(pwds) != 2 {
		panic("expect len==2")
	}
	fmt.Printf("elems:%+v\n", pwds)

	// query data in string
	// table type
	titles, data, err := database.QueryTable(mdb, "SELECT * FROM user LIMIT 10")
	if err != nil {
		panic(err)
	}
	fmt.Printf("table title:%+v\n", titles)
	fmt.Printf("table data: %+v\n", data)
	// map type
	mData, err := database.QueryMap(mdb, "SELECT * FROM user LIMIT 10")
	if err != nil {
		panic(err)
	}
	fmt.Printf("mData: %+v\n", mData)

	// execute for tx
	tx, err := mdb.Begin()
	if err != nil {
		panic(err)
	}
	txUsers := []TestingUser{
		{UserName: "t3", Passwd: "t3"},
		{UserName: "t4", Passwd: "t4"},
	}
	for _, u := range txUsers {
		if _, err := database.InsertStruct(tx, "user", &u); err != nil {
			fmt.Println(errors.As(err))
			database.Rollback(tx)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(errors.As(err))
		database.Rollback(tx)
		return
	}

	// excute for stmt
	stmt, err := mdb.Prepare("SELECT COUNT(*) FROM user WHERE username=?")
	count := 0
	if err := stmt.QueryRow("t3").Scan(&count); err != nil {
		panic(err)
	}
	if count != 1 {
		panic(errors.New("need count==1").As(count))
	}
}
