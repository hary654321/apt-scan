package scanner

import (
	"net"
	"sync"
	"zrWorker/core/cdn"
	"zrWorker/core/otherScanner"
	"zrWorker/core/slog"
	"zrWorker/lib/cache"
	"zrWorker/pkg/utils"
)

var CDNCheck = false

var DomainDatabase = sync.Map{}

type DomainClient struct {
	*client
	HandlerIsCDN  func(domain, CDNInfo string)
	HandlerRealIP func(domain string, ip net.IP)
	HandlerError  func(domain string, err error)
}

func NewDomainScanner(config *Config) *DomainClient {
	var client = &DomainClient{
		client:        newConfig(config, config.Threads),
		HandlerIsCDN:  func(domain, CDNInfo string) {},
		HandlerRealIP: func(domain string, ip net.IP) {},
		HandlerError:  func(domain string, err error) {},
	}
	client.pool.Interval = config.Interval
	client.pool.Function = func(in interface{}) {
		domain := in.(string)
		ip, err := cdn.Resolution(domain)
		if err != nil {
			client.HandlerError(domain, err)
			return
		}
		slog.Println(slog.DEBUG, domain, ip)
		client.HandlerRealIP(domain, net.ParseIP(ip))
		//将DNS解析结果存入数据库
		DomainDatabase.Store(domain, ip)
		if CDNCheck == false {
			slog.Println(slog.DEBUG, "HandlerRealIP", domain, ip)
			client.HandlerRealIP(domain, net.ParseIP(ip))
			return
		}

		if ok, result, _ := cdn.FindWithDomain(domain); ok {
			slog.Println(slog.DEBUG, domain, ip)
			client.HandlerIsCDN(domain, result)
			return
		}
		if ok, result, _ := cdn.FindWithIP(ip); ok {
			slog.Println(slog.DEBUG, domain, ip)
			client.HandlerIsCDN(domain, result)
			return
		}
	}
	return client
}

func (c *DomainClient) Push(domain string) {
	c.pool.Push(domain)
}

// 保存域名相关的信息
func (c *DomainClient) SaveDomainRootInfo(runTaskID, root string) {

	kkkk := root + "icp"
	if len(cache.Get(kkkk)) > 0 {
		slog.Println(slog.DEBUG, "根域名查询过Icp cert :", root)
	} else {
		m := make(map[string]interface{})
		m["Domain"] = root
		domainIp := utils.GetDomainIp(root)
		m["IP"] = domainIp

		//cerTInfo := cert.HttpGetCert("https://" + root)
		//if cerTInfo != nil {
		//	m["Cert_NotBefore"] = cerTInfo.NotBefore.String()
		//	m["Cert_NotAfter"] = cerTInfo.NotAfter.String()
		//	m["Cert_Subject"] = cerTInfo.Subject.String()
		//	m["Cert_Issuer"] = cerTInfo.Issuer.String()
		//	m["Cert_PublicKeyAlgorithm"] = cerTInfo.PublicKeyAlgorithm.String()
		//	m["Cert_SignatureAlgorithm"] = cerTInfo.SignatureAlgorithm.String()
		//	m["Cert_KeyUsage"] = utils.GetInterfaceToString(cerTInfo.KeyUsage)
		//}
		//slog.Println(slog.INFO, "SaveDomainInfo:", result)
		//m["Addr"] = result
		//obj, ok := icp.GetIcpInfo(root)
		//if ok {
		//	m["ICP_unitName"] = utils.GetInterfaceToString(obj["unitName"])
		//	m["ICP_updateRecordTime"] = utils.GetInterfaceToString(obj["updateRecordTime"])
		//	m["ICP_domain"] = utils.GetInterfaceToString(obj["domain"])
		//	m["ICP_natureName"] = utils.GetInterfaceToString(obj["natureName"])
		//	m["ICP_unitName"] = utils.GetInterfaceToString(obj["unitName"])
		//	m["ICP_serviceId"] = utils.GetInterfaceToString(obj["serviceId"])
		//	m["ICP_domainId"] = utils.GetInterfaceToString(obj["domainId"])
		//	m["ICP_serviceLicence"] = utils.GetInterfaceToString(obj["serviceLicence"])
		//	utils.WriteJsonAny(utils.GetLogPath(runTaskID, "icpInfo"), m)
		//}
		//whois 加入
		info, errw := otherScanner.GetWhoisInfo(root)
		if errw == nil {
			m["whois"] = info
			m["Port"] = 80
		}

		m["SaveDomainRootInfo"] = "SaveDomainRootInfo"
		utils.WriteJsonAny(utils.GetLogPath(runTaskID, "ipInfo"), m)

		time := utils.GetTime()
		cache.Set(kkkk, []byte(time))
	}
}

// 保存子域名信息
func (c *DomainClient) SaveSubDomain(runTaskID, root string) (subDomains []string) {

	kkkk := root + "subDomain"
	//slog.Println(slog.INFO, "SaveSubDomain：", root)
	if len(cache.Get(kkkk)) > 0 {
		slog.Println(slog.DEBUG, "子域名刚刚查询过:", root)
	} else {
		time := utils.GetTime()
		cache.Set(kkkk, []byte(time))
		subDomains = otherScanner.SaveSubDomain(runTaskID, root)
	}

	return
}
