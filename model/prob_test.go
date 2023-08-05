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

	resp, err := TcpSend("tcp", "192.168.56.132:6722", "0500000040000000a2-a4000000f09330819f300d06092a864886f70d010101050003818d0030818902818100e42b643814d3b9006fc4fbd6f50c5ace6aaedd2e5ea940ee8d8d1143c9a014d08ad7820c836f7bc355ba96db20f8d4830d52ed8373325e2b398b432e7cac71c4da3613c91a93791c285699fb38f405110ceee5922f2d515fb2af979df6fa324407489d55974338c33f38721d113d5b7dae7843f3b7913c29717ddbbb217db4430203010001", 10)
	log.Println("Scan:", resp, err)

}

func TestReg(t *testing.T) {
	str := "G1olang regular expressions example"
	match, err := regexp.MatchString(`^Golang`, str)
	fmt.Println("Match: ", match, " Error: ", err)
}
