package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	// CoreConf crocodile conf
	CoreConf *coreConf
)

// Init Config
func Init(conf string) {
	_, err := toml.DecodeFile(conf, &CoreConf)
	if err != nil {
		fmt.Printf("Err %v", err)
		os.Exit(1)
	}
}

type coreConf struct {
	ApiPort     int    `json:"apiPort"`
	HttpsServer bool   `json:"https"`
	BasicAuth   string `json:"basicAuth"`
	Version     string
}
