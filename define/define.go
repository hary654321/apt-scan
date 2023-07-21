package define

type Portres struct {
	CreateDate       string `json:"CreateDate"`
	CreateTime       string `json:"CreateTime"`
	DeviceType       string `json:"DeviceType"`
	Hostname         string `json:"Hostname"`
	IP               string `json:"IP"`
	Info             string `json:"Info"`
	MatchRegexString string `json:"MatchRegexString"`
	OperatingSystem  string `json:"OperatingSystem"`
	Port             string `json:"Port"`
	ProbeName        string `json:"ProbeName"`
	ProductName      string `json:"ProductName"`
	Response         string `json:"Response"`
	Service          string `json:"Service"`
	Version          string `json:"Version"`
	RunTaskID        string `json:"runTaskID"`
}
