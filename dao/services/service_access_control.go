package serviceDAO

import (
	"GoGateway/dao"
)

type ServiceAccessControl struct {
	ID                uint   `json:"id" gorm:"primary_key"`
	ServiceID         uint   `json:"service_id"`
	OpenAuth          int    `json:"open_auth" gorm:"type:tinyint"`
	BlackList         string `json:"black_list" gorm:"type:varchar(600)"`
	WhiteList         string `json:"white_list" gorm:"type:varchar(600)"`
	WhiteHostName     string `json:"white_host_name" gorm:"type:varchar(255)"`
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" gorm:"type:int"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"type:int"`
}

func (s *ServiceAccessControl) TableName() string {
	return "go_gateway_service_access_control"
}

// Handler Methods

type ServiceAccessControlHandler struct{}

func (sac *ServiceAccessControlHandler) Find(dt *ServiceAccessControl) (*ServiceAccessControl, error) {
	result := &ServiceAccessControl{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (sac *ServiceAccessControlHandler) Save(dt *ServiceAccessControl) error {
	return dao.DB.Save(dt).Error
}
