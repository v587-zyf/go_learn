package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type PluginMgr struct {
	// 维护所有插件
	plugins map[string]Registry
	lock    sync.Mutex
}

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

func RegisterPlugin(registry Registry) (err error) {
	return pluginMgr.registerPlugin(registry)
}

// 插件注册
func (p *PluginMgr) registerPlugin(registry Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.plugins[registry.Name()]; ok {
		err = errors.New("registry exist")
		return
	}

	p.plugins[registry.Name()] = registry

	return
}

// 初始化
func InitRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	return pluginMgr.initRegistry(ctx, name, opts...)
}

func (p *PluginMgr) initRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	plugin, ok := p.plugins[name]
	if !ok {
		err = fmt.Errorf("registry %s not exist", name)
		return
	}

	registry = plugin
	registry.Init(ctx, opts...)

	return
}
