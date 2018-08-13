package test

import (
	"testing"

	"github.com/gwaylib/conf/ini"
)

func TestEtc(t *testing.T) {
	cfg, err := ini.GetFile("./etc.cfg")
	if err != nil {
		t.Fatal(err)
	}
	sec := cfg.Section("test")
	if sec.Key("str").String() != "abc" {
		t.Fatal(sec.Key("str"))
	}
	if sec.Key("int").MustInt() != 1 {
		t.Fatal(sec.Key("int"))
	}
	if sec.Key("bool_true").MustBool() != true {
		t.Fatal(sec.Key("bool_true"))
	}
	if sec.Key("bool_false").MustBool() != false {
		t.Fatal(sec.Key("bool_false"))
	}
	if sec.Key("float").MustFloat64() != 3.20 {
		t.Fatal(sec.Key("float"))
	}
	for i, v := range sec.Key("slice").Ints("|") {
		if v-1 != i {
			t.Fatal(sec.Key("slice"))
		}
	}
}

func TestI18n(t *testing.T) {
	i18nDir := "./app.default."
	cfg := ini.NewIni(i18nDir)
	//	msg_en := cfg.Get("en").Section("error").Key("0").String()
	//	if msg_en != "zero" {
	//		t.Fatal(msg_en)
	//		return
	//	}
	//
	//	msg_zh_cn := cfg.Get("zh_cn").Section("error").Key("0").String()
	//	if msg_zh_cn != "零" {
	//		t.Fatal(msg_zh_cn)
	//		return
	//	}
	//
	//	msg_zh_tw := cfg.Get("zh_tw").Section("error").Key("0").String()
	//	if msg_zh_tw != "零" {
	//		t.Fatal(msg_zh_tw)
	//		return
	//	}

	msg_default := cfg.GetDefault("", "en").Section("error").Key("0").String()
	if msg_default != "zero" {
		t.Fatal(msg_default)
		return
	}
}
