package orm

import (
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
)

func WithDryRun() contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

func WithConfigPath(configPath string) contract.DBOption {

	return func(container framework.Container, config *contract.DBConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)

		//加载config配置路径
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}

}
