package etcd

import (
	"context"
	"encoding/json"
	"example/module/registry"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

type RegistryService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

type EtcdRegistry struct {
	options            *registry.Options
	client             *clientv3.Client
	serviceCh          chan *registry.Service
	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegistryService
}

const (
	MaxServiceNum = 8
)

var (
	etcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegistryService, MaxServiceNum),
	}
)

func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	etcdRegistry.value.Store(allServiceInfo)
	err := registry.RegisterPlugin(etcdRegistry)
	if err != nil {
		err = fmt.Errorf("init registerPlugin err:%v", err)
	}

	go etcdRegistry.run()
}

func (e *EtcdRegistry) run() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case service := <-e.serviceCh:
			registryService, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registryService.service.Nodes = append(registryService.service.Nodes, node)
				}
				registryService.registered = false
				break
			}
			registryService = &RegistryService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			e.syncServiceFromEtcd()
		default:
			e.registerOnKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (e *EtcdRegistry) syncServiceFromEtcd() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	ctx := context.TODO()
	allServiceInfo = e.value.Load().(*AllServiceInfo)
	for _, service := range allServiceInfo.serviceMap {
		key := e.servicePath(service.Name)
		rsp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfo.serviceMap[service.Name] = service
			continue
		}
		serviceNew := &registry.Service{
			Name: service.Name,
		}
		for _, kv := range rsp.Kvs {
			val := kv.Value
			var tmpService registry.Service
			err = json.Unmarshal(val, &tmpService)
			if err != nil {
				fmt.Printf("json unmarshal val:%v err:%v", val, err)
				return
			}
			for _, node := range tmpService.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}
		allServiceInfo.serviceMap[serviceNew.Name] = serviceNew
	}
	e.value.Store(allServiceInfo)
}

func (e *EtcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

func (e *EtcdRegistry) registerOnKeepAlive() {
	for _, registryService := range e.registryServiceMap {
		if registryService.registered {
			e.keepAlive(registryService)
			continue
		}
		err := e.registerService(registryService)
		if err != nil {
			fmt.Printf("registerService err:%v", err)
			continue
		}
	}
}

func (e *EtcdRegistry) keepAlive(registryService *RegistryService) {
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
	}
	return
}

func (e *EtcdRegistry) registerService(registryService *RegistryService) (err error) {
	resp, err := e.client.Grant(context.TODO(), e.options.Heartbeat)
	if err != nil {
		return
	}

	registryService.id = resp.ID
	for _, node := range registryService.service.Nodes {
		tmp := &registry.Service{
			Name: registryService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}
		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}

		key := e.serviceNodePath(tmp)
		fmt.Println("register key:", key)

		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}

		ch, err := e.client.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			continue
		}
		registryService.keepAliveCh = ch
		registryService.registered = true
	}

	return
}

func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIP := fmt.Sprintf("%s:%d", service.Nodes[0].Ip, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIP)
}

func (e *EtcdRegistry) Name() string {
	return "etcd"
}

// 初始化
func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})
	if err != nil {
		err = fmt.Errorf("init etcd err:%v", err)
		return
	}

	return
}

// 服务注册
func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceCh <- service:
	default:
		err = fmt.Errorf("serviceCh is full")
		return
	}

	return
}

// 服务反注册
func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {

	return
}

// 服务发现
func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	service, err = e.GetService(ctx, name)
	if err != nil {
		return
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	service, err = e.getServiceFromCache(ctx, name)
	if err != nil {
		return
	}

	key := e.servicePath(name)
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	service = &registry.Service{
		Name: name,
	}
	for _, kv := range resp.Kvs {
		val := kv.Value
		var tmpService registry.Service
		err = json.Unmarshal(val, &tmpService)
		if err != nil {
			return
		}
		for _, node := range tmpService.Nodes {
			service.Nodes = append(service.Nodes, node)
		}
	}
	allServiceInfoOld := e.value.Load().(*AllServiceInfo)
	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	for k, v := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[k] = v
	}
	allServiceInfoNew.serviceMap[name] = service
	e.value.Store(allServiceInfoNew)

	return
}

func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, err error) {
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	service, ok := allServiceInfo.serviceMap[name]
	if !ok {
		err = fmt.Errorf("service:%s not found", name)
		return
	}

	return
}
