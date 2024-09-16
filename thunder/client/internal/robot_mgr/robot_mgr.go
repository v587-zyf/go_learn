package robot_mgr

import (
	"client/internal/robot"
	"sync"
)

type RobotManager struct {
	options *RobotOption

	Robot map[uint64]*robot.Robot
	sync.RWMutex
}

func NewRobotManager() *RobotManager {
	r := &RobotManager{
		options: NewRobotOption(),
	}

	return r
}

func (r *RobotManager) Init(opts ...Option) error {
	for _, opt := range opts {
		opt(r.options)
	}

	return nil
}

func (r *RobotManager) Start() error {
	robotCfg := &robot.RobotConf{
		Https: r.options.https,
		Pem:   r.options.pem,
		Key:   r.options.key,
	}
	rbt := robot.NewRobot(robotCfg)
	rbt.Init()
	rbt.Run()

	//menu.Menus()

	return nil
}
func (r *RobotManager) Stop() {
	r.RLock()
	defer r.RUnlock()

	for _, rbt := range r.Robot {
		rbt.Stop(rbt.GetSession())
	}
}
