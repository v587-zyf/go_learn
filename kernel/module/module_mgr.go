package module

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"kernel/iface"
	"kernel/log"
	"sync"
)

type ModuleMgr struct {
	modules map[string]iface.IModule
	sync.RWMutex
}

func (mm *ModuleMgr) Add(m iface.IModule) {
	mm.Lock()
	defer mm.Unlock()

	if _, ok := mm.modules[m.Name()]; ok {
		log.Warn("module already exists", zap.String("name", m.Name()))
		return
	}

	mm.modules[m.Name()] = m
}

func (mm *ModuleMgr) Get(name string) iface.IModule {
	mm.RLock()
	defer mm.RUnlock()

	module, ok := mm.modules[name]
	if !ok {
		log.Warn("module not exists", zap.String("name", name))
		return nil
	}

	return module
}

func (mm *ModuleMgr) Del(name string) {
	mm.Lock()
	defer mm.Unlock()

	delete(mm.modules, name)
}

func (mm *ModuleMgr) Init(ctx context.Context, opts ...iface.Option) (err error) {
	mm.RLock()
	defer mm.RUnlock()

	if len(mm.modules) <= 0 {
		return fmt.Errorf("no module")
	}

	for _, module := range mm.modules {
		err = module.Init(ctx, opts...)
		if err != nil {
			return
		}
	}

	return nil
}

func (mm *ModuleMgr) Start() (err error) {
	mm.RLock()
	defer mm.RUnlock()

	if len(mm.modules) <= 0 {
		return fmt.Errorf("no module")
	}

	for _, module := range mm.modules {
		err = module.Start()
		if err != nil {
			return
		}
	}

	return nil
}

func (mm *ModuleMgr) Run() {
	mm.RLock()
	defer mm.RUnlock()

	for _, module := range mm.modules {
		module.Run()
	}
}

func (mm *ModuleMgr) Stop() {
	mm.RLock()
	defer mm.RUnlock()

	for _, module := range mm.modules {
		module.Run()
	}
}
