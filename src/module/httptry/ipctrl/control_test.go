package ipctrl

import (
	"testing"
)

// TODO:因已在真实环境测试,但仍需进行更多的攻击测试

func TestControl(t *testing.T) {
	ctrl := NewAccessController()
	if err := ctrl.SetWhiteList("./testdata/ip_access_white.json"); err != nil {
		t.Fatal(err)
	}
	if err := ctrl.SetBlackList("./testdata/ip_access_black.json"); err != nil {
		t.Fatal(err)
	}
	testIp := []string{"127.0.0.1"}
	if ip, ok := ctrl.InBlackList(testIp); !ok {
		t.Fatal("should be true")
	} else if ip != testIp[0] {
		t.Fatal(ip, testIp)
	}

	if ip, ok := ctrl.InWhiteList(testIp); !ok {
		t.Fatal("should be true")
	} else if ip != testIp[0] {
		t.Fatal(ip, testIp)
	}

	// 使用例子
}
