package serviceDAO

import "GoGateway/dao"

type ServiceLoadBalance struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	ServiceID     uint   `json:"service_id"`
	CheckMethod   int    `json:"check_method"`
	CheckTimeout  int    `json:"check_timeout"`
	CheckInterval int    `json:"check_interval"`
	RoundType     int    `json:"round_type"`
	IpList        string `json:"ip_list"`
	WeightList    string `json:"weight_list"`
	ForbidList    string `json:"forbid_list"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout"`
	UpstreamMaxIdle        int `json:"upstream_max_idle"`
}

func (s *ServiceLoadBalance) TableName() string {
	return "go_gateway_service_load_balance"
}

// Handler Methods

type ServiceLoadBalanceHandler struct{}

func (slb *ServiceLoadBalanceHandler) Find(dt *ServiceLoadBalance) (*ServiceLoadBalance, error) {
	result := &ServiceLoadBalance{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (slb *ServiceLoadBalanceHandler) Save(dt *ServiceLoadBalance) error {
	return dao.DB.Save(dt).Error
}
