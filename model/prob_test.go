package model

import (
	"log"
	"testing"
)

func Test_PayloadPreHandle(t *testing.T) {
	res := PayloadPreHandle("GET /sys.php HTTP/1.1\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8", "127.0.0.1")
	log.Println("Scan:", res)
}

func Test_tcp(t *testing.T) {

	resp, err := Conn("tcp", "192.168.56.132:6722", "0500000040000000a2", 10)
	log.Println("Scan:", resp, err)

}
