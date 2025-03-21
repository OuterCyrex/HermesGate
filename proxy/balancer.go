package proxy

import (
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"GoGateway/proxy/load_balance"
	"strings"
)

// ServiceBalancer 用于存储所有服务的负载均衡策略
type ServiceBalancer struct {
	ServiceBalanceMap map[string]*load_balance.LoadBalance
}

func NewServiceBalancer() *ServiceBalancer {
	return &ServiceBalancer{
		ServiceBalanceMap: make(map[string]*load_balance.LoadBalance),
	}
}

var ServiceBalanceHandler *ServiceBalancer

func init() {
	ServiceBalanceHandler = NewServiceBalancer()
}

// GetLoadBalance 用于通过 serviceDAO.ServiceDetail 获取对应的负载均衡策略
func (lbr *ServiceBalancer) GetLoadBalance(detail *serviceDAO.ServiceDetail) (load_balance.LoadBalance, error) {

	// 如果已经存在则不再进行初始化负载均衡
	if v, ok := lbr.ServiceBalanceMap[detail.Info.ServiceName]; ok {
		return *v, nil
	}

	ipList := strings.Split(detail.LoadBalance.IpList, ",")
	weightList := strings.Split(detail.LoadBalance.WeightList, ",")

	if len(ipList) != len(weightList) {
		return nil, status.Errorf(codes.InvalidParams, "IPList And WeightList length not match")
	}

	ipConfig := map[string]string{}

	// 将ip和权重信息装入负载均衡设置
	for index, ip := range ipList {
		ipConfig[ip] = weightList[index]
	}

	mConf := load_balance.NewLoadBalanceCheckConf("%s", ipConfig)

	lb := load_balance.LoadBalanceFactorWithConf(load_balance.LoadBalanceType(detail.LoadBalance.RoundType), mConf)

	lbr.ServiceBalanceMap[detail.Info.ServiceName] = &lb

	return lb, nil
}

func (lbr *ServiceBalancer) ReloadLoadBalance(detail *serviceDAO.ServiceDetail) error {
	ipList := strings.Split(detail.LoadBalance.IpList, ",")
	weightList := strings.Split(detail.LoadBalance.WeightList, ",")

	if len(ipList) != len(weightList) {
		return status.Errorf(codes.InvalidParams, "IPList And WeightList length not match")
	}

	ipConfig := map[string]string{}

	// 将ip和权重信息装入负载均衡设置
	for index, ip := range ipList {
		ipConfig[ip] = weightList[index]
	}

	mConf := load_balance.NewLoadBalanceCheckConf("%s", ipConfig)

	lb := load_balance.LoadBalanceFactorWithConf(load_balance.LoadBalanceType(detail.LoadBalance.RoundType), mConf)

	lbr.ServiceBalanceMap[detail.Info.ServiceName] = &lb

	return nil
}
