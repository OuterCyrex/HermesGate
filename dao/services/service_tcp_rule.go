package serviceDAO

import "GoGateway/dao"

type ServiceTcpRule struct {
	ID        uint `json:"id" gorm:"primary_key"`
	ServiceID uint `json:"serviceId"`
	Port      int  `json:"port"`
}

func (s *ServiceTcpRule) TableName() string {
	return "go_gateway_service_tcp_rule"
}

// Handler Methods

type ServiceTcpRuleHandler struct{}

func (str *ServiceTcpRuleHandler) Find(dt *ServiceTcpRule) (*ServiceTcpRule, error) {
	result := &ServiceTcpRule{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (str *ServiceTcpRuleHandler) Save(dt *ServiceTcpRule) error {
	return dao.DB.Save(dt).Error
}
