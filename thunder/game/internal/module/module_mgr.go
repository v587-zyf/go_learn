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

	mods := map[enum.ModuleName]iface.IModule{
		enum.G_M_ENTER:  NewEnterMgr(),
		enum.G_M_MAP:    NewMapMgr(),
		enum.G_M_USER:   NewUserMgr(),
		enum.G_M_CARD:   NewCardMgr(),
		enum.G_M_SHOP:   NewShopMgr(),
		enum.G_M_HASTEN: NewHastenMgr(),
		enum.G_M_INVITE: NewInviteMgr(),
		enum.G_M_RANK:   NewRankMgr(),
		enum.G_M_GUILD:  NewGuildMgr(),
	}
	for mName, mFn := range mods {
		m.AddModule(mName, mFn)
	}
	return m
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
