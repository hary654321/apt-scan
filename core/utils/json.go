package utils

import (
	"bytes"
	"encoding/json"
	"github.com/likexian/gokit/xjson"
	"os"
	"zrWorker/core/slog"
	"zrWorker/global"
)

func WriteJson(path string, data interface{}) {

	if global.ServerSetting.RunMode == "debug" {
		prettyData, err1 := xjson.PrettyDumps(data)
		if err1 != nil {
			slog.Println(slog.DEBUG, err1)
		}
		slog.Println(slog.DEBUG, path, prettyData)
	}

	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		slog.Println(slog.DEBUG, err)
	}

	f, err := os.OpenFile(path, os.O_CREATE+os.O_RDWR+os.O_APPEND, 0764)
	if err != nil {
		slog.Println(slog.DEBUG, err)
	}

	//jsonBuf := append([]byte(result),[]byte("\r\n")...)
	f.Write(buf.Bytes())

}

func WriteJsonAny(path string, m map[string]interface{}) {
	m["CreateDate"] = GetDate()
	m["CreateTime"] = GetTime()
	if global.ServerSetting.RunMode == "debug" {
		m["Body"] = ""
		m["FingerPrint"] = ""
		m["Response"] = ""
		m["FoundDomain"] = ""
	}

	WriteJson(path, m)
}

func WriteJsonString(path string, m map[string]string) {
	m["CreateDate"] = GetDate()
	m["CreateTime"] = GetTime()
	if global.ServerSetting.RunMode == "debug" {
		m["Body"] = ""
		m["FingerPrint"] = ""
		m["Response"] = ""
		m["FoundDomain"] = ""
	}

	WriteJson(path, m)
}
