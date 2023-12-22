package model

import (
	"fmt"
	"log"
	"regexp"
	"testing"
)

func Test_PayloadPreHandle(t *testing.T) {
	res := PayloadPreHandle("GET /sys.php HTTP/1.1\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8", "127.0.0.1")
	log.Println("Scan:", res)
}

func Test_tcp(t *testing.T) {

	resp, err := TcpSend("tcp", "127.0.0.1:445", "hi", 5)
	log.Println("Scan:", resp, err)

}

func TestReg(t *testing.T) {
	str := `</head>
	<body>404 not found,power by <a href=" //ehang.io/nps">nps</a>`
	match, err := regexp.MatchString(`*404 not found,power by*nps`, str)
	fmt.Println("Match: ", match, " Error: ", err)
}

// b5ce1a0a3349aeeb00000000a002faf036a70000020405b40402080a00851e480000000001030307
func Test_frp(t *testing.T) {

	resp, err := TcpSend("tcp", "192.168.56.132:6666", "7077640a", 8)
	log.Println("Scan:", resp, err)

}

func Test_NC(t *testing.T) {

	resp, err := TcpSend("UDP", "192.168.56.132:6666", "7077640a", 10)
	log.Println("Scan:", resp, err)

}

func Test_Msf(t *testing.T) {

	resp, err := TcpSend("tcp", "192.168.56.141:4444", "6c730a", 6)
	log.Println("Scan:", resp, err)

}
