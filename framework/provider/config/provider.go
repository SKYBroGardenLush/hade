package config

import (
	"github.com/SKYBroGardenLush/skycraper/framework"
	"github.com/SKYBroGardenLush/skycraper/framework/contract"
	"path/filepath"
)

type HadeConfigProvider struct{}

func (receiver *HadeConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeConfig
}

func (receiver *HadeConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)

	//配置文件夹地址
	configFolder := appService.ConfigFolder()
	env := envService.AppEnv()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

// Boot will called when the services instantiate
func (receiver *HadeConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the services instantiate when first make or register
func (receiver *HadeConfigProvider) IsDefer() bool {
	return false
}

// Name adefine the name for this services
func (receiver *HadeConfigProvider) Name() string {
	return contract.ConfigKey
}
