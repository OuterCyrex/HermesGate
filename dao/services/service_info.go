package serviceDAO

import (
	"GoGateway/biz/model/services"
	"GoGateway/dao"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"gorm.io/gorm"
)

type ServiceInfo struct {
	gorm.Model

	LoadType    int    `json:"load_type"`
	ServiceName string `json:"service_name" gorm:"type:varchar(256)"`
	ServiceDesc string `json:"service_description" gorm:"type:varchar(256)"`
}

func (s *ServiceInfo) TableName() string {
	return "go_gateway_service_info"
}

func (s *ServiceInfo) ToHttpResponse(addr string, nodes int) *services.ServiceListItemResponse {
	loadType := "http"

	switch s.LoadType {
	case 0:
		loadType = "http"
	case 1:
		loadType = "tcp"
	case 2:
		loadType = "grpc"
	}

	return &services.ServiceListItemResponse{
		Id:          int32(s.ID),
		ServiceName: s.ServiceName,
		ServiceDesc: s.ServiceDesc,
		LoadType:    loadType,
		ServiceAddr: addr,
		TotalNode:   int32(nodes),
		QPS:         0,
		Qpd:         0,
	}
}

// Repository Methods

type ServiceInfoRepository struct{}

func (s *ServiceInfoRepository) Delete(id uint) error {
	result := dao.DB.Delete(&ServiceInfo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "服务信息不存在")
	}
	return nil
}

func (s *ServiceInfoRepository) Find(dt *ServiceInfo) (*ServiceInfo, error) {
	result := &ServiceInfo{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
