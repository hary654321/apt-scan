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

	PortWin   = filepath.Join(GlobalPath, "scn/r/port/")
	PortLinux = "/opt/scn/r/port/"

	WebMgrWin   = filepath.Join(GlobalPath, "scn/r/webMgr/")
	WebMgrLinux = "/opt/scn/r/webMgr/"

	SrvIdentWin   = filepath.Join(GlobalPath, "scn/r/srvIdent/")
	SrvIdentLinux = "/opt/scn/r/srvIdent/"

	PicklePath      = "pickle"
	ServiceTypeNums = []string{"webDir", "passwd_crack", "sslCert", "probe", "port", "webMgr", "srvIdent"}
)

// GetPicklePaths 返回pickle路径
func GetPicklePaths() []string {
	paths := make([]string, 0)
	for _, path := range ServiceTypeNums {
		val := filepath.Join(GlobalPath, PicklePath, path)
		paths = append(paths, val)
	}
	return paths
}
