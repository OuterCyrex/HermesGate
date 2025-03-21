package proxy

import (
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/consts/codes"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/pkg/status"
	serviceSVC "GoGateway/svc/services"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
	"sync"
)

type ServiceManager struct {
	ServiceMap map[string]*serviceDAO.ServiceDetail
	Lock       *sync.RWMutex
	err        error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap: map[string]*serviceDAO.ServiceDetail{},
		Lock:       &sync.RWMutex{},
	}
}

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

// Load 用于将服务信息存入内存
func (s *ServiceManager) Load() error {
	svc := serviceSVC.ServiceInfoSvcLayer{}
	total, _, _ := svc.PageList("", 0, 0)
	_, serverList, err := svc.PageList("", 1, int32(total))
	if err != nil {
		s.err = err
	}

	for _, server := range serverList {
		detail, err := svc.ServiceDetail(&server)
		if err != nil {
			s.err = err
		}
		s.Lock.Lock()
		s.ServiceMap[server.ServiceName] = detail
		s.Lock.Unlock()
	}
	return s.err
}

func (s *ServiceManager) Reload(detail *serviceDAO.ServiceDetail) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.ServiceMap[detail.Info.ServiceName] = detail
}

// GetHttpDetail 用于根据用户访问的url获取对应的serviceDetail对象
func (s *ServiceManager) GetHttpDetail(c *app.RequestContext) (*serviceDAO.ServiceDetail, error) {
	host := string(c.Request.Host())

	idx := strings.Index(host, ":")
	if idx != -1 {
		host = host[:strings.Index(host, ":")]
	}

	path := string(c.Request.Path())

	for _, item := range s.ServiceMap {
		if item.Info.LoadType != serviceConsts.ServiceLoadTypeHTTP {
			continue
		}

		switch item.Http.RuleType {
		case serviceConsts.HTTPRuleTypeDomain:
			if item.Http.Rule == host {
				return item, nil
			} else {
				continue
			}
		case serviceConsts.HTTPRuleTypePrefixURL:
			if strings.HasPrefix(path, item.Http.Rule) {
				return item, nil
			} else {
				continue
			}
		default:
			continue
		}
	}

	return nil, status.Errorf(codes.NotFound, "cannot find certain http service")
}
