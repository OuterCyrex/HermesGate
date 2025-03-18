package serviceDAO

import (
	"GoGateway/biz/model/services"
	"GoGateway/dao"
	"GoGateway/pkg/consts/codes"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/pkg/status"
	"GoGateway/proxy/redis_counter"
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

	counter := redisCounter.ServiceFlowCountHandler.GetCounter(s.ServiceName)

	return &services.ServiceListItemResponse{
		Id:          int32(s.ID),
		ServiceName: s.ServiceName,
		ServiceDesc: s.ServiceDesc,
		LoadType:    loadType,
		ServiceAddr: addr,
		TotalNode:   int32(nodes),
		QPS:         counter.QPS,
		Qpd:         counter.TotalCount,
	}
}

// Repository Methods

type ServiceInfoRepository struct{}

func (s *ServiceInfoRepository) Delete(id uint) error {
	tx := dao.DB.Begin()
	var info ServiceInfo

	result := tx.Where(&ServiceInfo{Model: gorm.Model{ID: id}}).First(&info)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return status.Errorf(codes.NotFound, "服务信息不存在")
	}

	result = tx.Delete(&ServiceInfo{}, id)
	if result.Error != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	switch info.LoadType {
	case serviceConsts.ServiceLoadTypeHTTP:
		tx.Where(&ServiceHttpRule{ServiceID: id}).Delete(&ServiceHttpRule{})
		if tx.Error != nil {
			tx.Rollback()
			return status.Errorf(codes.InternalError, result.Error.Error())
		}
	case serviceConsts.ServiceLoadTypeGRPC:
		tx.Where(&ServiceGRPCRule{ServiceID: id}).Delete(&ServiceGRPCRule{})
		if tx.Error != nil {
			tx.Rollback()
			return status.Errorf(codes.InternalError, result.Error.Error())
		}

	case serviceConsts.ServiceLoadTypeTCP:
		tx.Where(&ServiceTcpRule{ServiceID: id}).Delete(&ServiceTcpRule{})
		if tx.Error != nil {
			tx.Rollback()
			return status.Errorf(codes.InternalError, result.Error.Error())
		}
	}

	tx.Where(ServiceAccessControl{ServiceID: id}).Delete(&ServiceAccessControl{})
	if tx.Error != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	tx.Where(ServiceLoadBalance{ServiceID: id}).Delete(&ServiceLoadBalance{})
	if tx.Error != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	tx.Commit()

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
