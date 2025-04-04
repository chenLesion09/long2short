package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(data []byte) string {
	encrypt := md5.New() // 创建加密器
	encrypt.Write(data)
	return hex.EncodeToString(encrypt.Sum(nil)) // 将32位16进制数转化为字符串并返回
}
