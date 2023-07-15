package udp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"testing"
	"time"
	"zrWorker/core/slog"
	"zrWorker/lib/ip/xdb"
)

func TestUdpServer(t *testing.T) {
	// 监听UDP服务
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})

	if err != nil {
		log.Fatal("Listen failed,", err)
		return
	}

	// 循环读取消息
	for {
		var data [1024]byte
		n, addr, err := udpConn.ReadFromUDP(data[:])
		if err != nil {
			log.Printf("Read from udp server:%s failed,err:%s", addr, err)
			break
		}
		go func() {
			// 返回数据
			fmt.Printf("Addr:%s,data:%v count:%d \n", addr, string(data[:n]), n)
			_, err := udpConn.WriteToUDP([]byte("msg recived."), addr)
			if err != nil {
				fmt.Println("write to udp server failed,err:", err)
			}
		}()
	}

}

// go test -v -run TestUdpClient port_test.go
func TestUdpClient(t *testing.T) {
	// 连接服务器
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9091,
	})

	if err != nil {
		log.Println("Connect to udp server failed,err:", err)
		return
	}

	// 发送数据
	_, err1 := conn.Write([]byte(fmt.Sprintf("udp testing:%v", "hi")))
	if err1 != nil {
		log.Println("Send data failed,err:", err)
		return
	}
	//接收数据
	result := make([]byte, 1024)

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		fmt.Printf("Listening...\n")
		//read from UDPConn here
		n, remoteAddr, err := conn.ReadFromUDP(result)
		if err != nil {
			log.Println("Read from udp server failed ,err:", err)
			return
		}
		fmt.Printf("Recived msg from %s, data:%s \n", remoteAddr, string(result[:n]))
	}
	fmt.Printf("Finished listening.\n")

}

func TestScanUDP(t *testing.T) {

	println(UpdStatus(net.IPv4(115, 45, 105, 240), 177))

}

// ListenICMP 拦截ICMP报文
func ListenICMP(ctx context.Context, address string, ch chan string) {

	go func() {
		for {
			TryUDP(address)
			time.Sleep(1 * time.Second)
		}

	}()

	netaddr, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, _ := net.ListenIP("ip4:icmp", netaddr)
	for {
		slog.Println(slog.DEBUG, "ListenICMP")
		buf := make([]byte, 1024)
		n, addr, _ := conn.ReadFrom(buf)
		msg, _ := icmp.ParseMessage(1, buf[0:n])
		fmt.Println(n, addr, msg.Type, msg.Code, msg.Checksum)

		ch <- addr.String()

		select {
		case <-ctx.Done():
			close(ch)
			log.Println(ctx.Err())
			return
		default:
			//fmt.Println(n, addr, msg.Type, msg.Code, msg.Checksum, string(marshal))
		}
	}

}

type UDPMessage struct {
	SrcPort  int
	DesPort  int
	Len      int
	CheckSum []byte
	Data     []byte
}

// ParseUDPMessage 解析UDP包
func ParseUDPMessage(b []byte) (*UDPMessage, error) {
	if len(b) < 8 {
		return nil, errors.New("invalid len")
	}
	m := &UDPMessage{}
	m.SrcPort = int(binary.BigEndian.Uint16(b[0:2]))
	m.DesPort = int(binary.BigEndian.Uint16(b[2:4]))
	m.Len = int(binary.BigEndian.Uint16(b[4:6]))
	m.CheckSum = b[6:8]
	m.Data = b[8:]
	return m, nil
}

// ParseICMPCode 解析ICMP类型和code
func ParseICMPCode(typeCode icmp.Type, code int) string {
	switch typeCode {
	case ipv4.ICMPTypeEchoReply:
		switch code {
		case 0:
			return "Echo Reply"
		}
	case ipv4.ICMPTypeDestinationUnreachable:
		switch code {
		case 0:
			return "Network Unreachable"
		case 1:
			return "Host Unreachable"
		case 2:
			return "Protocol Unreachable"
		case 3:
			return "Port Unreachable"
		}
	case ipv4.ICMPTypeEcho:
		return "Echo Request"
	}

	return fmt.Sprintf("未知CODE %s %d", typeCode, code)
}

var sendData = []byte("Hello Server")

// TryUDP 向目标端口发送UDP数据
func TryUDP(address string) error {
	slog.Println(slog.DEBUG, address)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		slog.Println(slog.DEBUG, "地址解析失败，err:", err)
	}
	socket, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		slog.Println(slog.DEBUG, "连接UDP服务器失败，err:", err)
	}
	defer socket.Close()

	_, err = socket.Write(sendData) // 发送数据
	if err != nil {
		slog.Println(slog.DEBUG, "发送数据失败，err:", err)
	}
	return nil
}

func TestIcmp(t *testing.T) {

	netaddr, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, _ := net.ListenIP("ip4:icmp", netaddr)
	for {
		buf := make([]byte, 1024)
		n, addr, _ := conn.ReadFrom(buf)
		msg, _ := icmp.ParseMessage(1, buf[0:n])
		fmt.Println(n, addr, msg.Type, msg.Code, msg.Checksum)
	}

}

func TestIc(t *testing.T) {

	a := xdb.GetIpUint32("0.0.0.0")
	b := 1674032790 - 1674032680
	println(a, b)
}
