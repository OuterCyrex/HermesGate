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
	ServiceBalanceList []*ServiceBalanceInfo
}

type ServiceBalanceInfo struct {
	ServiceName string
	LoadBalance *load_balance.LoadBalance
}

func NewServiceBalancer() *ServiceBalancer {
	return &ServiceBalancer{
		ServiceBalanceList: make([]*ServiceBalanceInfo, 0),
	}
}

var ServiceBalanceHandler *ServiceBalancer

func init() {
	ServiceBalanceHandler = NewServiceBalancer()
}

// GetLoadBalance 用于通过 serviceDAO.ServiceDetail 获取对应的负载均衡策略
func (lbr *ServiceBalancer) GetLoadBalance(detail *serviceDAO.ServiceDetail) (load_balance.LoadBalance, error) {

	// 如果已经存在则不再进行初始化负载均衡
	for _, lbrItem := range lbr.ServiceBalanceList {
		if lbrItem.ServiceName == detail.Info.ServiceName {
			return *lbrItem.LoadBalance, nil
		}
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

	lbItem := &ServiceBalanceInfo{
		LoadBalance: &lb,
		ServiceName: detail.Info.ServiceName,
	}
	lbr.ServiceBalanceList = append(lbr.ServiceBalanceList, lbItem)

	return lb, nil
}
