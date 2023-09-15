package main

import (
	"bytes"
	"encoding/base64"
)

func main() {
	DecodeString()
}

func DecodeString() {
	decodedAuth, err := base64.StdEncoding.DecodeString("enc6Y2g=")
	if err != nil {
		// 处理解码错误
		return
	}
	// 按照冒号分割解码后的字符串，得到用户名和密码
	parts := bytes.Split(decodedAuth, []byte(":"))
	if len(parts) != 2 {
		// 处理格式错误
		return
	}
	username := string(parts[0])
	password := string(parts[1])

	println(username, password)
}
