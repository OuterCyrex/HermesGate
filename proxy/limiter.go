package proxy

import (
	"GoGateway/pkg/consts/rater"
	"golang.org/x/time/rate"
	"sync"
)

type limiter struct {
	serviceName string
	limiter     *rate.Limiter
}

type ServiceLimiter struct {
	limiterMap map[string]*limiter
	lock       sync.Locker
}

var ServiceLimitHandler *ServiceLimiter

func NewServiceLimiter() *ServiceLimiter {
	return &ServiceLimiter{
		limiterMap: make(map[string]*limiter),
		lock:       &sync.Mutex{},
	}
}

func init() {
	ServiceLimitHandler = NewServiceLimiter()
}

func (sl *ServiceLimiter) getLimiter(serviceName string, qps float64) *rate.Limiter {
	if v, exists := sl.limiterMap[serviceName]; exists {
		return v.limiter
	}

	sl.lock.Lock()
	defer sl.lock.Unlock()
	lim := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	sl.limiterMap[serviceName] = &limiter{
		serviceName: serviceName,
		limiter:     lim,
	}

	return lim
}

func (sl *ServiceLimiter) GetServerLimiter(serviceName string, qps float64) *rate.Limiter {
	return sl.getLimiter(raterConsts.RateLimiterKey+"_"+serviceName, qps)
}

func (sl *ServiceLimiter) GetClientLimiter(serviceName string, ip string, qps float64) *rate.Limiter {
	return sl.getLimiter(raterConsts.RateLimiterKey+"_"+ip+"_"+serviceName, qps)
}
