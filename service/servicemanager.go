package service

import (
	"sync"
)

type IServiceManager interface {
	Setup(s IService) bool
	Init(logger ILogger) bool
	Start() bool
	CreateServiceID() int
}

type CServiceManager struct {
	genserviceid    int
	localserviceMap map[string]IService
	logger          ILogger
}

func (slf *CServiceManager) Setup(s IService) bool {
	slf.localserviceMap[s.GetServiceName()] = s
	return true
}

func (slf *CServiceManager) FindService(serviceName string) IService {
	service, ok := slf.localserviceMap[serviceName]
	if ok {
		return service
	}

	return nil
}

type FetchService func(s IService) error

func (slf *CServiceManager) FetchService(s FetchService) IService {
	for _, se := range slf.localserviceMap {
		s(se)
	}

	return nil
}

func (slf *CServiceManager) Init(logger ILogger, exit chan bool, pwaitGroup *sync.WaitGroup) bool {
	slf.logger = logger
	for _, s := range slf.localserviceMap {
		(s.(IModule)).InitModule(exit, pwaitGroup)
	}

	return true
}

func (slf *CServiceManager) Start() bool {
	for _, s := range slf.localserviceMap {
		go (s.(IModule)).RunModule(s.(IModule))
	}

	return true
}

func (slf *CServiceManager) GenServiceID() int {
	slf.genserviceid += 1
	return slf.genserviceid
}

func (slf *CServiceManager) GetLogger() ILogger {
	return slf.logger
}

var self *CServiceManager

func InstanceServiceMgr() *CServiceManager {
	if self == nil {
		self = new(CServiceManager)
		self.localserviceMap = make(map[string]IService)
		return self
	}

	return self
}

func GetLogger() ILogger {
	return InstanceServiceMgr().GetLogger()
}
