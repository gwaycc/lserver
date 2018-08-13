package cms

import (
	"testing"
)

var testCmsPriv = []struct {
	priv   CmsPriv
	path   string
	result bool
}{
	{CmsPriv{}, "/", true}, // 所有都可以访问根节点
	{CmsPriv{}, "", true},  // 所有都可以访问根节点
	{CmsPriv{}, "test", false},
	{CmsPriv{}, "//", true},
	{CmsPriv{[]string{"user"}}, "/user", true},
	{CmsPriv{[]string{"user", "pwd"}}, "/user", true},
	{CmsPriv{[]string{"user"}}, "/user/pwd", false},
	{CmsPriv{[]string{"*"}}, "/user/pwd", false},
	{CmsPriv{[]string{"*", "*"}}, "/user/pwd", true},
}

func TestCmsPriv(t *testing.T) {
	for i, c := range testCmsPriv {
		out := c.priv.CheckSum(c.path)
		if out != c.result {
			t.Fatal(i, out, c)
		}
	}
}
