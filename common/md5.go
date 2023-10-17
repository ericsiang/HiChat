package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)


// Md5encoder 加密后返回小写值
func Md5encoder(code string) string {
	m := md5.New()
	io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

// Md5StrToUpper 加密后返回大写
func Md5StrToUpper(code string) string {
	return strings.ToUpper(Md5encoder(code))
}

func SaltPassword(password string, salt string) string {
	saltPwd := fmt.Sprintf("%s$%s", Md5encoder(password), salt)
	return saltPwd
}

func CheckPassWord(checkPassword string, salt string ,orignalPassword string) bool {
	return orignalPassword == SaltPassword(checkPassword, salt)
}
