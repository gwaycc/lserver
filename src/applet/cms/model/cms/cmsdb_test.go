package cms

import (
	"fmt"
	"testing"
	"time"
)

func TestCmsUser(t *testing.T) {
	db := NewCmsDB()
	defer db.Close()

	now := time.Now()
	username := fmt.Sprint(now.UnixNano())
	if len(username) > 32 {
		username = username[len(username)-32:]
	}
	if err := db.CreateUser(username, "123456", "test"); err != nil {
		t.Fatal(err)
	}
	if err := db.ResetPwd(username, "654321"); err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateUserStatus(username, 2); err != nil {
		t.Fatal(err)
	}
	u, err := db.GetUser(username, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !u.CheckSumPasswd("654321") {
		t.Fatal(u)
	}
	if u.NickName != "test" {
		t.Fatal(u.NickName)
	}
}

func TestUserPriv(t *testing.T) {
	db := NewCmsDB()
	defer db.Close()

	now := time.Now()
	username := fmt.Sprint(now.UnixNano())
	if len(username) > 32 {
		username = username[len(username)-32:]
	}

	path := "/test/" + username
	if err := db.CreateMenu("test."+username, "testing"); err != nil {
		t.Fatal(err)
	}
	priv, err := db.GetPriv(username)
	if err != nil {
		t.Fatal(err)
	}
	if priv.CheckSum(path) {
		t.Fatal(path)
	}

	if err := db.AddPriv(username, "test"); err != nil {
		t.Fatal(err)
	}
	priv, err = db.GetPriv(username)
	if err != nil {
		t.Fatal(err)
	}
	if !priv.CheckSum(path) {
		t.Fatal(priv, path)
	}
}
func TestBindUserPriv(t *testing.T) {
	db := NewCmsDB()
	defer db.Close()
	now := time.Now()
	username := fmt.Sprint(now.UnixNano())
	if len(username) > 32 {
		username = username[len(username)-32:]
	}
	if err := db.BindPriv(username, "管理员"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminPriv(t *testing.T) {
	db := NewCmsDB()
	defer db.Close()

	priv, err := db.GetPriv("admin")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(priv.ToJson())
}

func TestAdminUser(t *testing.T) {
	db := NewCmsDB()
	defer db.Close()

	u, err := db.GetUser("admin", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
}

func TestPutLog(t *testing.T) {
	db := NewCmsDB()
	defer db.Close()

	if err := db.PutLog("admin", -1, "testing", "go testing"); err != nil {
		t.Fatal(err)
	}
}
