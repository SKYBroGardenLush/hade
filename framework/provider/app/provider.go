package app

import (
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
)

type HadeAppProvider struct {
	BaseFolder string
}

func (h *HadeAppProvider) Name() string {
	return contract.AppKey
}

func (h *HadeAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeApp
}

// Params 获取初始化参数
func (h *HadeAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

func (h *HadeAppProvider) Boot(container framework.Container) error {
	return nil
}

func (h *HadeAppProvider) IsDefer() bool {
	return false
}
