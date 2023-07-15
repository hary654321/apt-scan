package utils

import (
	"testing"
	"zrWorker/lib/ip/xdb"
	ip2domian "zrWorker/lib/ip2domain"
)

func TestIcp(t *testing.T) {

	a := GetStrACount("a", "abca")
	PrinfI("a", a)
}

func TestStr(t *testing.T) {

	res := SubStrBefore("127.0.0.1:80", ":")
	println(res)

	resa := SubStrAfter("127.0.0.1:80", ":")
	println(resa)
}

func TestInArr(t *testing.T) {

	println(In_array("google.com", ip2domian.Top100DomainArr))
}

func TestWrtie(t *testing.T) {
	Write("a", "14.133.22.242")

	ipLast := Read("a")

	b := xdb.GetIpUint32(ipLast)

	println(b)
}

func TestIsInnerIp(t *testing.T) {
	res := IsInnerIp("172.16.0.1")
	println(res)
}
