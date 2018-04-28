package g

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/file"
	"sync"
)

type GoRadiusConfig struct {
	Debug        bool   `json:"debug"`
	RadiusDb     string `json:"radiusDb"`
	FireSystemDb string `json:"fireSystemDb"`
	AuthListen   string `json:"authListen"`
	AcctListen   string `json:"acctListen"`
	SharedKey    string `json:"sharedKey"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type GlobalConfig struct {
	GoRadius *GoRadiusConfig `json:"goRadius"`
	Http     *HttpConfig     `json:"http"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		fmt.Println("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		fmt.Println("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		fmt.Println("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		fmt.Println("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	fmt.Println("read config file:", cfg, "successfully")

}
