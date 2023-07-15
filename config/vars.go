package config

import (
	"os"
	"path/filepath"
)

var (
	GlobalPath, _ = os.Getwd()

	WebDirWin   = filepath.Join(GlobalPath, "scn/r/webDirSearch/")
	WebDirLinux = "/opt/scn/r/webDirSearch"

	PasswdCrackWin   = filepath.Join(GlobalPath, "scn/r/passwordCrack/")
	PasswdCrackLinux = "/opt/scn/r/passwordCrack/"

	SslCertWin   = filepath.Join(GlobalPath, "scn/r/certificate/")
	SslCertLinux = "/opt/scn/r/certificate/"

	ProbeWin   = filepath.Join(GlobalPath, "scn/r/probe/")
	ProbeLinux = "/opt/scn/r/probe/"

	WebMgrWin   = filepath.Join(GlobalPath, "scn/r/webMgr/")
	WebMgrLinux = "/opt/scn/r/webMgr/"

	SrvIdentWin   = filepath.Join(GlobalPath, "scn/r/srvIdent/")
	SrvIdentLinux = "/opt/scn/r/srvIdent/"

	PicklePath      = "pickle"
	ServiceTypeNums = []string{"webDir", "passwd_crack", "sslCert", "probe", "webMgr", "srvIdent"}

	NacosConfigList = make(map[string][2]string)
)

func init() {
	NacosConfigList["user-agent"] = [2]string{"iastool-user-agent", "USER_AGENT"}

	NacosConfigList["psd-user"] = [2]string{"iastool-crack-user", "PASSWD_CRACK"}
	NacosConfigList["psd-passwd"] = [2]string{"iastool-crack-passwd", "PASSWD_CRACK"}

	NacosConfigList["wdb-mini"] = [2]string{"iastool-web-dict-mini", "WEB_DICT"}
	NacosConfigList["wdb-normal"] = [2]string{"iastool-web-dict-normal", "WEB_DICT"}
	NacosConfigList["wdb-big"] = [2]string{"iastool-web-dict-big", "WEB_DICT"}

	NacosConfigList["iasTool"] = [2]string{"ias-tool", "DEFAULT_GROUP"}
}

//GetPicklePaths 返回pickle路径
func GetPicklePaths() []string {
	paths := make([]string, 0)
	for _, path := range ServiceTypeNums {
		val := filepath.Join(GlobalPath, PicklePath, path)
		paths = append(paths, val)
	}
	return paths
}
