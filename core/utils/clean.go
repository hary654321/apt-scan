package utils

import "strings"

func Dealdata(m map[string]string) map[string]string {
	if m["ProductName"] == "UOS" {
		m["OperatingSystem"] = "统信UOS"
	}

	if strings.Contains(m["Response"], "UOS_Desktop_104") {
		m["OperatingSystem"] = "统信UOS"
	}

	return m
}
