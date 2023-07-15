package initialize

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"ias_tool_v2/config"
	"ias_tool_v2/https"
)

type Config struct {
	ApiPort     int  `json:"apiPort"`
	HttpsServer bool `json:"https"`
}

var GConfig *Config

//GetIasConfigFromNacos 获取全局配置
func GetIasConfigFromNacos() (string, error) {

	midConfig := config.NacosConfigList["iasTool"]

	dataId, groupId := midConfig[0], midConfig[1]

	configClient, err := config.GetNacosClient()

	if err != nil {
		return "", err
	}

	if data, err := configClient.GetConfig(vo.ConfigParam{DataId: dataId, Group: groupId}); err != nil {
		return "", err
	} else {
		return data, nil
	}
}

// InitConfig 加载配置
func InitConfig() (err error) {
	var (
		conf Config
	)

	https.PemInit()

	data, err := GetIasConfigFromNacos()
	if err = json.Unmarshal([]byte(data), &conf); err != nil {
		return
	}
	GConfig = &conf
	return
}
