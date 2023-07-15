package scanner

import (
	"net"
	"zrWorker/core/slog"
	"zrWorker/lib/cache"
	"zrWorker/lib/ip"
	ip2domian "zrWorker/lib/ip2domain"
	"zrWorker/lib/ping"
	"zrWorker/lib/tcpping"
	"zrWorker/pkg/utils"
)

type IPClient struct {
	*client
	HandlerAlive func(addr net.IP)
	HandlerDie   func(addr net.IP)
	HandlerError func(addr net.IP, err error)
}

func NewIPScanner(config *Config) *IPClient {
	var client = &IPClient{
		client:       newConfig(config, config.Threads),
		HandlerAlive: func(addr net.IP) {},
		HandlerDie:   func(addr net.IP) {},
		HandlerError: func(addr net.IP, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		ip := in.(net.IP)
		if ping.Pinger(ip.String(), 1) == nil {
			client.HandlerAlive(ip)
			return
		}
		if err := tcpping.PingPorts(ip.String(), config.Timeout); err == nil {
			client.HandlerAlive(ip)
			return
		}
		client.HandlerDie(ip)
	}
	return client
}

func (c *IPClient) Push(ips ...net.IP) {
	//这里是1-255的IP
	//println(ips)
	//slog.Println(slog.DEBUG, "pushIp：", ips)
	for _, ip := range ips {
		if len(cache.Get(ip.String())) > 0 {
			slog.Println(slog.WARN, "ip刚刚扫描过:", ip.String())
		} else {
			c.pool.Push(ip)
			time := utils.GetTime()
			//otherScanner.SaveSubDomain(domain)
			cache.Set(ip.String(), []byte(time))
		}
	}
}

func (c *IPClient) SaveIpInfo(runTaskID, ipstr string, cha bool) (IpDomain []string) {

	m := make(map[string]string)
	m["IP"] = ipstr
	slog.Println(slog.DEBUG, "SearchIpAddr")
	iPInfo := ip.SearchIpAddr(ipstr)
	m["IP_country"] = iPInfo.Country
	m["IP_province"] = iPInfo.Province
	m["IP_city"] = iPInfo.City
	m["IP_isp"] = iPInfo.Isp
	//m["OperatingSystem"] = otherScanner.GetOpInfo(ipstr)

	if !utils.IsInnerIp(ipstr) && cha {
		slog.Println(slog.WARN, "IpDomain : ", ipstr)
		//AsnInfo := asn.GetAsnRt(ipstr)
		//if AsnInfo != nil {
		//	m["Asn"] = utils.GetInterfaceToString(AsnInfo["asn"])
		//	m["Asn_start_ip"] = utils.GetInterfaceToString(AsnInfo["prefix"])
		//	m["Asn_end_ip"] = utils.GetInterfaceToString(AsnInfo["prefix"])
		//	m["Asn_info"] = utils.GetInterfaceToString(AsnInfo["info"])
		//	m["Asn_type"] = utils.GetInterfaceToString(AsnInfo["info"])
		//	m["Asn_domain"] = utils.GetInterfaceToString(AsnInfo["info"])
		//}
		IpDomain = saveIpDomains(ipstr)
		m["IpDomain"] = utils.ArrayToString(IpDomain)
	}

	utils.WriteJsonString(utils.GetLogPath(runTaskID, "ipInfo"), m)

	return
}

func saveIpDomains(ipstr string) []string {
	domains := ip2domian.GETIpDomains(ipstr)

	var erJ []string
	for _, d := range domains {
		domain := d["domain"]
		slog.Println(slog.DEBUG, "ip-域名 ", ipstr, domain)
		if utils.GetStrACount(".", domain) != 1 {
			erJ = append(erJ, domain)
		}
	}

	return erJ
}
