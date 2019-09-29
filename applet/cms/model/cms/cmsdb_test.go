package cms

import (
	"fmt"
	"testing"
	"time"
)

func TestCmsUser(t *testing.T) {
	db := NewCmsDB()

	now := time.Now()
	username := fmt.Sprint(now.UnixNano())
	if len(username) > 32 {
		username = username[len(username)-32:]
	}
	if err := db.CreateUser(&CmsUser{
		UserName: username,
		Passwd:   "123456",
		NickName: "test",
	}); err != nil {
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

func TestGroupPriv(t *testing.T) {
	db := NewCmsDB()

	now := time.Now()
	gid := int(now.Unix())
	gidStr := fmt.Sprint(gid)

	path := "/test/" + gidStr
	if err := db.CreateMenu("test."+gidStr, "testing"); err != nil {
		t.Fatal(err)
	}
	priv, err := db.GetPriv(gid)
	if err != nil {
		t.Fatal(err)
	}
	if priv.CheckSum(path) {
		t.Fatal(path)
	}

	if err := db.AddPriv(gid, "test."+gidStr); err != nil {
		t.Fatal(err)
	}
	priv, err = db.GetPriv(gid)
	if err != nil {
		t.Fatal(err)
	}
	if !priv.CheckSum(path) {
		t.Fatal(priv, path)
	}
}
func TestBindGroupPriv(t *testing.T) {
	db := NewCmsDB()
	now := time.Now()
	gid := int(now.Unix())
	if err := db.BindPriv(gid, "管理员"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminPriv(t *testing.T) {
	db := NewCmsDB()

	priv, err := db.GetPriv(0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(priv.ToJson())
}

func TestAdminUser(t *testing.T) {
	db := NewCmsDB()

	u, err := db.GetUser("admin", 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
}

func TestPutLog(t *testing.T) {
	db := NewCmsDB()

	if err := db.PutLog("admin", "testing", "testing", "go testing"); err != nil {
		t.Fatal(err)
	}
}
