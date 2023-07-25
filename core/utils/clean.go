package utils

import "strings"

var OsMap map[string]string = map[string]string{
	"UOS":               "统信UOS",
	"UOS_Desktop_104":   "统信UOS",
	"(Uos)":             "统信UOS",
	"AQUOS":             "统信UOS",
	"(Ubuntu)":          "Ubuntu",
	"Set-Cookie: AIROS": "Ubiquiti AirOS路由器",
	"ciscoSystems":      "Cisco SNMP service",
	"openeuler":         "openEuler OS",
}

func Dealdata(m map[string]string) map[string]string {

	for k, v := range OsMap {
		if strings.Contains(m["Response"], k) {
			m["OperatingSystem"] = v
		}
	}

	return m
}
