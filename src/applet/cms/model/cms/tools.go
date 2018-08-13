package cms

import (
	"github.com/jameskeane/bcrypt"
)

// 生成密码
func CreatePwd(in string) (pwd, salt string) {
	salt, _ = bcrypt.Salt(10)
	pwd, _ = bcrypt.Hash(in, salt)
	return
}
