package cms

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestPasswd(t *testing.T) {
	pwdHash := fmt.Sprintf("%x", sha256.Sum256([]byte("LogAdmin123")))
	fmt.Println(pwdHash)
	pwd, _ := CreatePwd(pwdHash)
	fmt.Println(pwd)
}
