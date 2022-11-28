package orm

import (
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
)

// GormProvider Gorm 的服务提供者
type GormProvider struct {
}

func (p *GormProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeGorm
}

func (p *GormProvider) Boot(container framework.Container) error {
	return nil
}

func (p *GormProvider) IsDefer() bool {
	return true
}

func (p *GormProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (p *GormProvider) Name() string {
	return contract.ORMKey
}
