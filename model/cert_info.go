package model

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TcpConfig struct {
	Host     string
	Port     uint64
	Timeout  uint64
	ReadByte int
	Data     string
}

type TcpConn struct {
	Err string
	Req net.Conn
}

type TcpResult struct {
	Err  string
	Code int
	Data []byte
}

type TlsConn struct {
	Err string
	Req *tls.Conn
}

const (
	IsNotAlive = -1 + iota
	IsNotTls
	IsTls
)

//DoTlsRequest 开始发起tls请求
func (conn *TlsConn) DoTlsRequest(tc *TcpConfig) *TlsResult {
	if conn.Err != "" || conn.Req == nil {
		return &TlsResult{Err: conn.Err, Code: -1, Data: []byte{}}
	}
	var buf []byte
	var n int
	var err error
	tr := &TlsResult{}
	defer conn.Req.Close()
	state := conn.Req.ConnectionState()
	if len(state.PeerCertificates) > 0 {
		curCert := state.PeerCertificates[0]
		tr.UrlPort = tc.Host + ":" + strconv.Itoa(int(tc.Port))

		tr.Cert = CertData{
			Subject:       curCert.Subject.String(),
			CertSubjectCN: curCert.Subject.CommonName,
			CertSubjectC:  strings.Join(curCert.Subject.Country, ","),
			CertSubjectO:  strings.Join(curCert.Subject.Organization, ","),

			CertIssuer:   curCert.Issuer.String(),
			CertIssuerCN: curCert.Issuer.CommonName,
			CertIssuerC:  strings.Join(curCert.Issuer.Country, ","),
			CertIssuerO:  strings.Join(curCert.Issuer.Organization, ","),

			ValidNotbefore: curCert.NotBefore.Format("2006-01-02 15:04:05"),
			ValidNotafter:  curCert.NotAfter.Format("2006-01-02 15:04:05"),
			SerialNumber:   hex.EncodeToString(curCert.SerialNumber.Bytes()),
			CertFile:       base64.StdEncoding.EncodeToString(curCert.Raw),
		}

		if len(curCert.DNSNames) == 0 {
			tr.Cert.SubjectAlternativeName = []string{}
		} else {
			tr.Cert.SubjectAlternativeName = curCert.DNSNames
		}

		var tbuf bytes.Buffer
		for _, b := range sha1.Sum(curCert.Raw) {
			_, _ = fmt.Fprintf(&tbuf, "%02x", b)
		}
		tr.Cert.Thumbprint = tbuf.String()
	}
	if tc.Data != "" {
		_, err = conn.Req.Write([]byte(tc.Data))
		if err != nil {
			tr.Code = -2
			tr.Err = BuildErr(false, err, "DoTlsRequest write")
		}
	}
	if tc.ReadByte < 0 {
		buf, err = ioutil.ReadAll(conn.Req)
	} else if tc.ReadByte > 0 {
		buf = make([]byte, n)
		_, err = conn.Req.Read(buf)
	}
	if err != nil {
		tr.Code = -2
		tr.Err = BuildErr(false, err, "DoTcpRequest read")
		tr.Data = []byte{}
		return tr
	}
	tr.Data = buf
	return tr
}

type CertData struct {
	Subject                string   `json:"cert_subject"`     //证书Subject字段
	SubjectAlternativeName []string `json:"dns_names"`        //Subject Alternative Names字段
	ValidNotbefore         string   `json:"valid_from"`       //证书有效期从
	ValidNotafter          string   `json:"valid_to"`         //证书有效期止
	SerialNumber           string   `json:"cert_serialno"`    //证书序列号
	Thumbprint             string   `json:"cert_fingerprint"` //证书指纹 --唯一键 SHA1格式
	CertFile               string   `json:"cert_base64"`      //证书Base64
	CertSubjectCN          string   `json:"cert_subject_cn"`  //证书Subject-CN字段
	CertSubjectC           string   `json:"cert_subject_c"`   //证书Subject-C字段
	CertSubjectO           string   `json:"cert_subject_o"`   //证书Subject-O字段
	CertIssuer             string   `json:"cert_issuer"`      //证书Issuer-CN字段
	CertIssuerCN           string   `json:"cert_issuer_cn"`   //证书Issuer-CN字段
	CertIssuerC            string   `json:"cert_issuer_c"`    //证书Issuer-C字段
	CertIssuerO            string   `json:"cert_issuer_o"`    //证书Issuer-O字段
}

type TlsResult struct {
	Err     string
	Code    int
	Data    []byte
	Cert    CertData
	UrlPort string
}

var TlsRespRe = regexp.MustCompile("\x16\x03[\x00-\x03]")

func CheckIsTlsAndParseCert(addr string) (int, CertData) {
	var err error
	isTls := IsNotAlive
	certData := CertData{}
	isCertOk := false
	timeout := time.Duration(5) * time.Second
	// Send Client Hello
	hs := TlsClientHelloPackage
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		defer conn.Close()
		err = conn.SetWriteDeadline(time.Now().Add(timeout))
		if err == nil {
			if n, err := conn.Write([]byte(hs)); err == nil {
				_ = n
				recv := make([]byte, 8192)
				err = conn.SetReadDeadline(time.Now().Add(timeout))
				if err == nil {
					if n, err = conn.Read(recv[:]); err == nil {
						matches := TlsRespRe.FindAllSubmatchIndex(recv, -1)
						if len(matches) > 0 && matches[0][0] == 0 {
							isTls = IsTls
							if len(matches) > 1 { // only Certificate
								for _, match := range matches[1:] {
									certData, isCertOk = GetCertByHandshakeRecv(recv, match[0])
									if isCertOk {
										break
									}
								}
							} else { // Multiple Handshake Message
								certData, isCertOk = GetCertByHandshakeRecv(recv, 0)
							}
						} else {
							isTls = IsNotTls
						}
					}
				}
			}
		}
	} else {
		isTls = IsNotAlive
	}
	return isTls, certData
}

func CheckIsTlsFullAndParseCert(addr string) (int, CertData) {
	var err error
	isTls := IsNotAlive
	certData := CertData{}

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: time.Duration(1) * time.Second},
		"tcp",
		addr,
		&tls.Config{InsecureSkipVerify: true})
	if err == nil {
		isTls = IsTls
		defer conn.Close()
		state := conn.ConnectionState()
		if len(state.PeerCertificates) > 0 {
			cert := state.PeerCertificates[0]
			certData = GetCertFromX509(cert)
		}
	} else {
		log.Println(err)
		if !strings.HasPrefix(err.Error(), "dial tcp") {
			isTls = IsNotTls
		}
	}
	return isTls, certData
}

func GetCertByHandshakeRecv(recv []byte, index int) (CertData, bool) {
	var certData CertData
	var ok bool
	// recv: 以\x16\x03开头的TLS握手连接响应
	// index: 需要解析证书部分的起始字节位置
	if index == 0 { // multiple handshake message
		index += 2 + 1 + 2 // \x16\x03 + TLS版本 + 包长度
		index += 2         //ServerHello: sign + '\x00'
		shLen := int(binary.BigEndian.Uint16(recv[index : index+2]))
		if shLen > 8192 {
			return certData, false
		}
		index += 2 + shLen //server hello len + server hello
	} else { // only certificate
		index += 2 + 1 + 2 // \x16\x03 + TLS版本 + 包长度
	}
	if index > 8192-1 {
		return certData, false
	}
	sign := recv[index : index+1]
	if string(sign) == "\x0b" { // Certificate标志
		index += 1 + 3 + 3 + 1                                         // 标志位 + \x00 + 包长度 + \x00 + 包长度 + \x00
		certLen := int(binary.BigEndian.Uint16(recv[index : index+2])) //解析证书长度
		if certLen > 8192 {
			return certData, false
		}
		index += 2                              //证书长度
		certByte := recv[index : index+certLen] //证书字节
		certData = ParseCertByte(certByte)
		ok = true
	}
	return certData, ok
}

func ParseCertByte(certByte []byte) CertData {
	var certData CertData
	cert, err := x509.ParseCertificate(certByte)
	if err == nil {
		certData = GetCertFromX509(cert)
	}
	return certData
}

func GetCertFromX509(cert *x509.Certificate) CertData {
	var certData CertData
	certData = CertData{
		Subject:                cert.Subject.String(),
		CertSubjectCN:          cert.Subject.CommonName,
		CertSubjectC:           strings.Join(cert.Subject.Country, ","),
		CertSubjectO:           strings.Join(cert.Subject.Organization, ","),
		SubjectAlternativeName: cert.DNSNames,

		CertIssuer:   cert.Issuer.String(),
		CertIssuerCN: cert.Issuer.CommonName,
		CertIssuerC:  strings.Join(cert.Issuer.Country, ","),
		CertIssuerO:  strings.Join(cert.Issuer.Organization, ","),

		ValidNotbefore: cert.NotBefore.Format("2006-01-02 15:04:05"),
		ValidNotafter:  cert.NotAfter.Format("2006-01-02 15:04:05"),
		SerialNumber:   hex.EncodeToString(cert.SerialNumber.Bytes()),
		CertFile:       base64.StdEncoding.EncodeToString(cert.Raw),
	}
	var cBuf bytes.Buffer
	for _, b := range sha1.Sum(cert.Raw) {
		_, _ = fmt.Fprintf(&cBuf, "%02x", b)
	}
	certData.Thumbprint = cBuf.String()
	return certData
}

//NewTcpConfig 生成一个tcpConfig
func NewTcpConfig(host string, port uint64, timeout uint64, readByte int, data string) *TcpConfig {
	tc := &TcpConfig{
		Host:     host,
		Port:     port,
		Timeout:  timeout,
		ReadByte: readByte,
		Data:     data,
	}
	if tc.Timeout < 1 {
		tc.Timeout = 1
	}
	if tc.Host == "" {
		BuildErr(true, errors.New("empty host, change to 127.0.0.1"), "NewHttpConfig")
		tc.Host = "127.0.0.1"
	}
	if tc.Port < 1 || tc.Port > 65535 {
		BuildErr(true, errors.New("port error, change to 80"), "NewHttpConfig")
		tc.Port = 80
	}
	if tc.ReadByte < -1 {
		tc.ReadByte = -1
	}
	return tc
}

//NewTlsConn 新建新的tls连接
func NewTlsConn(tc *TcpConfig) (conn TlsConn, err error) {
	retryTimes := 3
	for i := 0; i < retryTimes; i++ {
		_conn, err := tls.DialWithDialer(
			&net.Dialer{Timeout: time.Duration(tc.Timeout) * time.Second},
			"tcp",
			fmt.Sprintf("%s:%d", tc.Host, tc.Port),
			&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			if err.Error() == "dial tcp: i/o timeout" {
				continue
			}
			log.Println("ERROR", err.Error())
			conn.Err = BuildErr(false, err, "NewTlsConn")
			return conn, err
		}

		if _conn != nil {
			conn.Req = _conn
			_ = conn.Req.SetDeadline(time.Now().Add(time.Duration(tc.Timeout) * time.Second))
			err = nil
		}
		return conn, nil
	}
	return conn, err
}
