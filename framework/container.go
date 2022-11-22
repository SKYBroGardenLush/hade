package framework

import (
	"errors"
	"sync"
)

// Container 是一个容器，提供绑定和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在会进行替换操作，返回error
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)

	// MustMake 根据关键字获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic
	//所以在使用这个接口的时候请保证服务容器已经为这个关键字绑定的服务提供者了
	MustMake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单列模式的
	//他时根据服务提供者注册的启动函数和传递的params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

func NewHeadContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instance:  map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

type HadeContainer struct {
	Container // 强制要求 HadeContainer 实现 Container 接口

	//存储注册服务的提供者，key为字符串凭证
	providers map[string]ServiceProvider
	//存储具体实例，key为字符串凭证
	instance map[string]interface{}
	// lock 用于锁住对容器变更操作
	lock sync.RWMutex
}

func (h *HadeContainer) Bind(provider ServiceProvider) error {
	h.lock.Lock()

	key := provider.Name()
	h.providers[key] = provider

	h.lock.Unlock()

	if provider.IsDefer() == false {
		if err := provider.Boot(h); err != nil {
			return err
		}

		//实例化方法
		params := provider.Params(h)

		method := provider.Register(h)

		instance, err := method(params...)
		if err != nil {
			return err
		}
		h.instance[key] = instance

	}
	return nil

}

//查询注册了的服务的提供者
func (h *HadeContainer) findServiceProvider(key string) ServiceProvider {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if provider, ok := h.providers[key]; ok {
		return provider
	} else {
		return nil
	}

}

func (h *HadeContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(h); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(h)
	}
	method := sp.Register(h)
	ins, err := method(params...)
	if err != nil {
		return nil, err
	}
	return ins, err
}

func (h *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {

	h.lock.RLock()
	defer h.lock.RUnlock()
	//查询是否已经注册了这个服务的提供者，如果没有注册，则返回错误

	sp := h.findServiceProvider(key)

	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return h.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器已经实例化了，那么就直接使用容器中的实例
	if ins, ok := h.instance[key]; ok {
		return ins, nil
	}

	// 如果容器还未实例化则进行一次实例化

	ins, err := h.newInstance(sp, params)
	if err != nil {
		return nil, err
	}

	h.instance[key] = ins
	return ins, err
}

func (h *HadeContainer) Make(key string) (interface{}, error) {
	return h.make(key, nil, false)
}

func (h *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return h.make(key, params, true)
}

func (h *HadeContainer) MustMake(key string) interface{} {
	serv, err := h.make(key, nil, false)

	if err != nil {
		panic(err)
	}
	return serv
}

// NameList 列出容器中所有服务提供者的字符串凭证
func (h *HadeContainer) NameList() []string {
	ret := []string{}
	for _, provider := range h.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}

func (h *HadeContainer) IsBind(key string) bool {
	if _, ok := h.providers[key]; ok {
		return true
	}
	return false
}
