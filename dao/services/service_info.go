package serviceDAO

import (
	"GoGateway/biz/model/services"
	"GoGateway/dao"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"errors"
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

func (s *ServiceInfo) ToHttpResponse(addr string, nodes int) *services.ServiceInfoResponse {
	loadType := "http"

	switch s.LoadType {
	case 0:
		loadType = "http"
	case 1:
		loadType = "rpc"
	case 2:
		loadType = "grpc"
	}

	return &services.ServiceInfoResponse{
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

// Handler Methods

type ServiceInfoHandler struct{}

func (s *ServiceInfoHandler) Delete(id uint) error {
	result := dao.DB.Delete(&ServiceInfo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "服务信息不存在")
	}
	return nil
}

func (s *ServiceInfoHandler) ServiceDetail(dt *ServiceInfo) (*ServiceDetail, error) {
	httpHandler := &ServiceHttpRuleHandler{}
	httpRule, err := httpHandler.Find(&ServiceHttpRule{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	grpcHandler := &ServiceGRPCRuleHandler{}
	grpcRule, err := grpcHandler.Find(&ServiceGRPCRule{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	tcpHandler := &ServiceTcpRuleHandler{}
	tcpRule, err := tcpHandler.Find(&ServiceTcpRule{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	loadHandler := &ServiceLoadBalanceHandler{}
	loadRule, err := loadHandler.Find(&ServiceLoadBalance{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	accessHandler := &ServiceAccessControlHandler{}
	accessRule, err := accessHandler.Find(&ServiceAccessControl{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	return &ServiceDetail{
		Info:          dt,
		Http:          httpRule,
		Tcp:           tcpRule,
		Grpc:          grpcRule,
		LoadBalance:   loadRule,
		AccessControl: accessRule,
	}, nil
}

func (s *ServiceInfoHandler) PageList(info string, pageNum int32, pageSize int32) (int64, []ServiceInfo, error) {
	var serviceList []ServiceInfo

	query := dao.DB.Model(ServiceInfo{})

	if info != "" {
		query = query.Where("service_name like ? OR service_desc like ?", "%"+info+"%", "%"+info+"%")
	}

	var count int64

	query.Count(&count)

	result := query.Scopes(dao.Paginate(pageNum, pageSize)).Find(&serviceList)
	if result.Error != nil {
		return 0, nil, status.Errorf(codes.InternalError, "Error while paginating service: %v", result.Error)
	}

	return count, serviceList, nil
}

func (s *ServiceInfoHandler) Find(dt *ServiceInfo) (*ServiceInfo, error) {
	result := &ServiceInfo{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Interfaces for integrated operations

func (s *ServiceInfoHandler) NewHTTPService(req services.ServiceAddHTTPRequest) (*ServiceInfo, error) {
	tx := dao.DB.Begin()

	serviceInfo := &ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
	}

	result := tx.Create(&serviceInfo)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InternalError, result.Error.Error())
	}

	httpHandler := &ServiceHttpRuleHandler{}
	if err := httpHandler.Save(&ServiceHttpRule{
		ServiceID:      serviceInfo.ID,
		RuleType:       int(req.RuleType),
		Rule:           req.Rule,
		NeedHttps:      int(req.NeedHTTP),
		NeedWebsocket:  int(req.NeedWebsocket),
		NeedStripUri:   int(req.NeedStripUri),
		UrlRewrite:     req.UrlRewrite,
		HeaderTransfer: req.HeaderTransfer,
	}); err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	accessHandler := &ServiceAccessControlHandler{}
	if err := accessHandler.Save(&ServiceAccessControl{
		ServiceID:         serviceInfo.ID,
		OpenAuth:          int(req.OpenAuth),
		BlackList:         req.BlackList,
		WhiteList:         req.WhiteList,
		ClientIPFlowLimit: int(req.ClientIPFlowLimit),
		ServiceFlowLimit:  int(req.ServiceFlowLimit),
	}); err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	loadBalanceHandler := &ServiceLoadBalanceHandler{}
	if err := loadBalanceHandler.Save(&ServiceLoadBalance{
		ServiceID:              serviceInfo.ID,
		RoundType:              int(req.RoundType),
		IpList:                 req.IpList,
		WeightList:             req.WeightList,
		UpstreamConnectTimeout: int(req.UpstreamConnectTimeout),
		UpstreamHeaderTimeout:  int(req.UpstreamHeaderTimeout),
		UpstreamIdleTimeout:    int(req.UpstreamIdleTimeout),
		UpstreamMaxIdle:        int(req.UpstreamMaxIdle),
	}); err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	return nil, nil
}
