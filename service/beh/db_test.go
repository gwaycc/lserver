package main

import (
	"testing"
	"time"

	"lserver/module/gouuid"
)

func TestSqlite3(t *testing.T) {
	if err := InsertLog(&DbTable{
		md5:      gouuid.New(),
		platform: "testing",
		ip:       "127.0.0.1",
		date:     time.Now(),
		level:    0,
		logger:   "testing",
		msg:      "testing",
	}); err != nil {
		t.Fatal(err)
	}
}
