package conf

import (
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
)

var (
	instance *MainConfig
	once     sync.Once
)

func GetConfig() *MainConfig {
	once.Do(func() {
		instance = &MainConfig{}
		instance.InitConfig()
	})
	return instance
}

func (c *MainConfig) InitConfig() {
	YAMLFile := "conf/config.yaml"
	v := viper.New()
	v.SetConfigFile(YAMLFile)
	if err := v.ReadInConfig(); err != nil {
		hlog.Fatalf("failed to read config: %v", err)
	}

	if err := v.Unmarshal(c); err != nil {
		hlog.Fatalf("failed to unmarshal config: %v", err)
	}
}
