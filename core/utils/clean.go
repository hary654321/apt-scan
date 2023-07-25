package utils

import "strings"

var OsMap map[string]string = map[string]string{
	"UOS_Desktop_104":   "统信UOS",
	"(Uos)":             "统信UOS",
	"AQUOS":             "统信UOS",
	"(Ubuntu)":          "Ubuntu",
	"Set-Cookie: AIROS": "Ubiquiti AirOS路由器",
	"ciscoSystems":      "Cisco SNMP service",
	"openeuler":         "openEuler OS",
	"Deepin/":           "深度操作系统",
	"(Deepin)":          "深度操作系统",
	"(NeoKylin":         "中标麒麟操作系统",
	"(Kylin)":           "银河麒麟Kylin OS",
	"Asianux":           "红旗Asianux操作系统",
}

func Dealdata(m map[string]string) map[string]string {

	for k, v := range OsMap {
		if strings.Contains(m["Response"], k) {
			m["OperatingSystem"] = v
		}
	}

	return m
}
