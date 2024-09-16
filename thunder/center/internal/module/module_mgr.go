package module

import (
	"comm/t_enum"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type ClientModuleMgr struct {
	module.DefModule

	options *ModuleMgrOption

	childM map[enum.ModuleName]iface.IModule
}

var cmMgr *ClientModuleMgr

func init() {
	cmMgr = NewClientModuleMgr()
}

func GetClientModuleMgr() *ClientModuleMgr {
	return cmMgr
}
func GetClientModuleMgrOptions() *ModuleMgrOption {
	return cmMgr.options
}

func Init(opts ...Option) error {
	for _, opt := range opts {
		opt(cmMgr.options)
	}

	return nil
}
func Start() { cmMgr.Start() }
func Run()   { cmMgr.Run() }
func Stop()  { cmMgr.Stop() }

func NewClientModuleMgr() *ClientModuleMgr {
	m := &ClientModuleMgr{
		childM:  make(map[enum.ModuleName]iface.IModule),
		options: NewModuleMgrOption(),
	}

	m.AddModule(enum.C_M_STRENGTH, NewStrengthMgr())
	m.AddModule(enum.C_M_GOLD, NewGoldMgr())
	return m
}

func InitOption(opts ...Option) error {
	for _, opt := range opts {
		opt(cmMgr.options)
	}

	return nil
}

func (m *ClientModuleMgr) AddModule(mn enum.ModuleName, module iface.IModule) {
	m.childM[mn] = module
}
func (m *ClientModuleMgr) GetModule(mn enum.ModuleName) iface.IModule {
	return m.childM[mn]
}

func (m *ClientModuleMgr) Init() {
	for _, iModule := range m.childM {
		clsName := fmt.Sprintf("%T", iModule)
		dotIndex := strings.Index(clsName, ".") + 1
		//log.Info(clsName[dotIndex:len(clsName)] + " Init")
		if err := iModule.Init(m.GetCtx()); err != nil {
			log.Error("module init err", zap.String("name", clsName[dotIndex:len(clsName)]))
		}
	}
}

func (m *ClientModuleMgr) Start() error {
	for _, iModule := range m.childM {
		err := iModule.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *ClientModuleMgr) Run() {
	for _, iModule := range m.childM {
		iModule.Run()
	}
}

func (m *ClientModuleMgr) Stop() {
	var wg sync.WaitGroup
	for _, iModule := range m.childM {
		wg.Add(1)
		go func(module iface.IModule) {
			module.Stop()
			wg.Done()
		}(iModule)
	}
	wg.Wait()
}
