package utils

import "strings"

var InnerPre = []string{"127.0.0.1", "10", "20", "172", "192"}

func IsInnerIp(ip string) bool {
	res := false
	for _, v := range InnerPre {
		if strings.HasPrefix(ip, v) {
			res = true
		}
	}
	return res
}
