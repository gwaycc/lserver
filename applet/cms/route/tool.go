package route

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gwaylib/errors"
	"github.com/xxtea/xxtea-go/xxtea"
)

var (
	XXTeaKey = []byte("nKaLi3hhb%5Fd9")
)

// XXTea加密
//
// 参数
// v -- 需要加密的对象
//
// 返回
// 返回经过json序列化、XXTea加密、base64.URLEncoding编码后的字符串
func XXTeaEncode(v interface{}) string {
	jData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(xxtea.Encrypt(jData, XXTeaKey))
}

// XXTea解密
//
// 参数
// src -- 加密的字符源
// v -- 需要json反序化出来的对象, 请参考json.Unmarshal
func XXTeaDecode(src string, v interface{}) error {
	xxteaData, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return errors.As(err)
	}
	if err := json.Unmarshal(xxtea.Decrypt(xxteaData, XXTeaKey), v); err != nil {
		return errors.As(err)
	}
	return nil
}
