package model

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"ias_tool_v2/config"
	"ias_tool_v2/core/slog"
	"io"
	"math"
	"net"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

var TlsClientHelloPackage = "\x16\x03\x00\x00\x69\x01\x00\x00\x65\x03\x03U\x1c\xa7\xe4random1random2random3random4\x00\x00\x0c\x00/\x00\x0a\x00\x13\x009\x00\x04\x00\xff\x01\x00\x00\x30\x00\x0d\x00,\x00*\x00\x01\x00\x03\x00\x02\x06\x01\x06\x03\x06\x02\x02\x01\x02\x03\x02\x02\x03\x01\x03\x03\x03\x02\x04\x01\x04\x03\x04\x02\x01\x01\x01\x03\x01\x02\x05\x01\x05\x03\x05\x02"

const IsTLS = 1
const IsNotTLS = 0

// LoadServiceResMap 根据传入的service_type 加载文件总路径
func LoadServiceResMap() (ServiceResMap map[string]map[string]string) {

	ServiceResMap = make(map[string]map[string]string)
	WebDir := make(map[string]string)
	PasswdCrack := make(map[string]string)
	SslCert := make(map[string]string)
	probe := make(map[string]string)
	webMgr := make(map[string]string)
	srvIdent := make(map[string]string)

	WebDir["windows"] = config.WebDirWin
	WebDir["linux"] = config.WebDirLinux

	PasswdCrack["windows"] = config.PasswdCrackWin
	PasswdCrack["linux"] = config.PasswdCrackLinux

	SslCert["windows"] = config.SslCertWin
	SslCert["linux"] = config.SslCertLinux

	probe["windows"] = config.ProbeWin
	probe["linux"] = config.ProbeLinux

	webMgr["windows"] = config.WebMgrWin
	webMgr["linux"] = config.WebMgrLinux

	srvIdent["windows"] = config.SrvIdentWin
	srvIdent["linux"] = config.SrvIdentLinux

	ServiceResMap["webDir"] = WebDir
	ServiceResMap["psd"] = PasswdCrack
	ServiceResMap["sslCert"] = SslCert
	ServiceResMap["probe"] = probe
	ServiceResMap["webMgr"] = webMgr
	ServiceResMap["srvIdent"] = srvIdent

	return ServiceResMap
}

// PathExist 判断文件是否存在
func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// GetResultPath 获取存放结果的总目录
func GetResultPath(serviceType string) string {
	sysName := runtime.GOOS
	ServiceResMap := LoadServiceResMap()
	return ServiceResMap[serviceType][sysName]
}

// InitResPath 初始化结果文件夹路径。没有则执行创建逻辑
func InitResPath() {
	for _, serviceName := range config.ServiceTypeNums {
		resPath := GetResultPath(serviceName)

		slog.Println(slog.DEBUG, "InitResPath", "resPath", resPath)
		if !PathExist(resPath) {
			_ = os.MkdirAll(resPath, os.ModePerm)
		}
	}
}

// InitPickle 初始化task pickle文件夹。没有则执行创建逻辑。
func InitPickle() {
	for _, path := range config.GetPicklePaths() {
		if !PathExist(path) {
			_ = os.MkdirAll(path, os.ModePerm)
		}
	}
}

// MinInt 由于math.min只支持float 先这么封装下吧
func MinInt(a, b int) (c int) {
	return int(math.Min(float64(a), float64(b)))
}

// String2Bytes 字符串转bytes
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Byte2GzipBase64Encoding(buf []byte) (base64Encoded string, err error) {
	var buffer bytes.Buffer
	zipBuf := gzip.NewWriter(&buffer)
	if _, err = zipBuf.Write(buf); err != nil {
		slog.Println(slog.DEBUG, "ERROR", "zip error", err.Error())
		return "", err
	}
	if err = zipBuf.Close(); err != nil {
		slog.Println(slog.DEBUG, "ERROR", "zip error", err.Error())
		return "", err
	}
	base64Encoded = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return base64Encoded, nil
}

// String2GzipBase64Encoding  字符串类型gzip压缩并进行base64加密
func String2GzipBase64Encoding(str string) (base64Encoded string, err error) {
	var buffer bytes.Buffer
	midByte := String2Bytes(str)
	zipBuf := gzip.NewWriter(&buffer)
	if _, err = zipBuf.Write(midByte); err != nil {
		slog.Println(slog.DEBUG, "ERROR", "zip error", err.Error())
		return "", err
	}
	if err = zipBuf.Close(); err != nil {
		slog.Println(slog.DEBUG, "ERROR", "zip error", err.Error())
		return "", err
	}
	base64Encoded = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return base64Encoded, nil
}

// BuildErr 封装error
func BuildErr(show bool, err error, msgs ...string) string {
	var _err string
	for i, msg := range msgs {
		_err += msg
		if i < len(msgs)-1 {
			_err += " "
		} else {
			_err += ":"
		}
	}
	_err += err.Error()
	if show {
		slog.Println(slog.DEBUG, "ERROR", _err)
	}
	return _err
}

func TcpSend(protocol, addr, payload string, timeout int) (string, error) {
	var err error
	var conn net.Conn

	conn, err = net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)

	if err != nil {
		slog.Println(slog.DEBUG, "conn:", err)
		return "", err
	}
	defer conn.Close()

	if len(payload) > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

		if strings.Contains(payload, "-") {
			parr := strings.Split(payload, "-")
			for _, p := range parr {

				pbyte, _ := hex.DecodeString(p)
				_, err = conn.Write(pbyte)
			}
		} else {
			pbyte, _ := hex.DecodeString(payload)

			_, err = conn.Write(pbyte)
		}
		if err != nil {
			slog.Println(slog.DEBUG, "write:", err)
			return "", err
		}
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	cr := bufio.NewReader(conn)
	var buf bytes.Buffer
	for {
		tmp, err := cr.ReadBytes('\n')
		if err == nil || err == io.EOF {
			ioErr := err
			_, err = buf.Write(tmp)
			if err != nil {
				slog.Println(slog.DEBUG, "read write2buf:", err)
			}
			if ioErr == io.EOF {
				break
			}
		} else {
			slog.Println(slog.DEBUG, "read in for:", err)
			break
		}
	}
	slog.Println(slog.DEBUG, "dump", hex.Dump(buf.Bytes()))

	return hex.EncodeToString(buf.Bytes()), nil
}

func HttpSend(protocol, addr, payload string, timeout int) (string, error) {
	var err error
	var conn net.Conn

	if protocol == "tls" {
		conn, err = tls.DialWithDialer(
			&net.Dialer{Timeout: time.Duration(5) * time.Second},
			"tcp",
			addr,
			&tls.Config{InsecureSkipVerify: true},
		)
	} else {
		conn, err = net.DialTimeout("tcp", addr, time.Duration(5)*time.Second)
	}
	if err != nil {
		slog.Println(slog.DEBUG, "conn:", err)
		return "", err
	}
	defer conn.Close()
	if len(payload) > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(5) * time.Second))
		_, err = conn.Write([]byte(payload))
		if err != nil {
			slog.Println(slog.DEBUG, "write:", err)
			return "", err
		}
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
	cr := bufio.NewReader(conn)
	var buf bytes.Buffer
	for {
		tmp, err := cr.ReadBytes('\n')
		if err == nil || err == io.EOF {
			ioErr := err
			_, err = buf.Write(tmp)
			if err != nil {
				slog.Println(slog.DEBUG, "read write2buf:", err)
			}
			if ioErr == io.EOF {
				break
			}
		} else {
			slog.Println(slog.DEBUG, "read in for:", err)
			break
		}
	}
	slog.Println(slog.DEBUG, buf.String())
	return buf.String(), nil
}

func CheckIsTls(addr string) int {
	var err error
	timeout := time.Duration(1) * time.Second

	hs := TlsClientHelloPackage
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		err = conn.SetWriteDeadline(time.Now().Add(timeout))
		if err == nil {
			if n, err := conn.Write([]byte(hs)); err == nil {
				recv := make([]byte, 8)
				err = conn.SetReadDeadline(time.Now().Add(timeout))
				if err == nil {
					if n, err = conn.Read(recv[:]); err == nil {
						re := regexp.MustCompile("^\x16\x03[\x00-\x03]")
						if isMatch := re.MatchString(string(recv[:n])); isMatch {
							return IsTLS
						}
					}
				}
			}
		}
		return IsNotTLS
	} else {
		return IsNotAlive
	}
}
