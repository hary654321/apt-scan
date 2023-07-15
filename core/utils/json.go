package utils

import (
	"bytes"
	"encoding/json"
	"ias_tool_v2/core/slog"

	"os"
)

func WriteJson(path string, data interface{}) {

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

	WriteJson(path, m)
}

func WriteJsonString(path string, m map[string]string) {
	m["CreateDate"] = GetDate()
	m["CreateTime"] = GetTime()

	WriteJson(path, m)
}
