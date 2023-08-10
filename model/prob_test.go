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

	resp, err := TcpSend("tcp", "192.168.56.132:8024", "545354-07000000302e32362e3130", 5)
	log.Println("Scan:", resp, err)

}

func TestReg(t *testing.T) {
	str := `</head>
	<body>404 not found,power by <a href=" //ehang.io/nps">nps</a>`
	match, err := regexp.MatchString(`*404 not found,power by*nps`, str)
	fmt.Println("Match: ", match, " Error: ", err)
}

func Test_frp(t *testing.T) {

	resp, err := TcpSend("tcp", "192.168.56.139:7000", "9a1fdc41ffedd3648a038f23c6e08774-00000000000000010000000b-5096d592ef438333af4cee-000000000000000100000010-94b2af7a914cd002ae67f519288a2fe0000000000000000100000043", 5)
	log.Println("Scan:", resp, err)

}
