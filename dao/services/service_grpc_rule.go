package serviceDAO

import "GoGateway/dao"

type ServiceGRPCRule struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	ServiceID   uint   `json:"serviceId"`
	Port        int    `json:"port"`
	HeaderTrans string `json:"header_trans"`
}

func (sgr *ServiceGRPCRule) TableName() string {
	return "go_gateway_service_grpc_rule"
}

// Handler Methods

type ServiceGRPCRuleHandler struct{}

func (sgr *ServiceGRPCRuleHandler) Find(dt *ServiceGRPCRule) (*ServiceGRPCRule, error) {
	result := &ServiceGRPCRule{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sgr *ServiceGRPCRuleHandler) Save(dt *ServiceGRPCRule) error {
	return dao.DB.Save(dt).Error
}
