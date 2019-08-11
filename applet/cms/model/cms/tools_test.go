package cms

import (
	"fmt"
	"testing"
)

func TestPasswd(t *testing.T) {
	pwd, _ := CreatePwd("Nd278@yS0587")
	fmt.Println(pwd)
}
