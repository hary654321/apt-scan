package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strings"
)

// GetNacosClient 客户端
func GetNacosClient() (configClient config_client.IConfigClient, err error) {
	sc := []constant.ServerConfig{
		{
			IpAddr: "129.226.189.251",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "a9ddf94e-704e-4240-983c-be755189ac26", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "error",
	}

	return clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
}

func GetConfigFromNacos(configKey string) ([]string, error) {
	config := NacosConfigList[configKey]
	dataId, groupId := config[0], config[1]

	configClient, err := GetNacosClient()

	if err != nil {
		return []string{}, err
	}

	if data, err := configClient.GetConfig(vo.ConfigParam{DataId: dataId, Group: groupId}); err != nil {
		return []string{}, err
	} else {
		configs := strings.Split(data, "\r\n")
		return configs, nil
	}
}
