package serviceSVC

import (
	"GoGateway/biz/model/services"
	"GoGateway/dao"
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"errors"
	"gorm.io/gorm"
)

type ServiceInfoSvcLayer struct{}

func (s *ServiceInfoSvcLayer) ServiceDetail(dt *serviceDAO.ServiceInfo) (*serviceDAO.ServiceDetail, error) {
	httpRepository := &serviceDAO.ServiceHttpRuleRepository{}
	httpRule, err := httpRepository.Find(&serviceDAO.ServiceHttpRule{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	grpcRepository := &serviceDAO.ServiceGRPCRuleRepository{}
	grpcRule, err := grpcRepository.Find(&serviceDAO.ServiceGRPCRule{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	tcpRepository := &serviceDAO.ServiceTcpRuleRepository{}
	tcpRule, err := tcpRepository.Find(&serviceDAO.ServiceTcpRule{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	loadRepository := &serviceDAO.ServiceLoadBalanceRepository{}
	loadRule, err := loadRepository.Find(&serviceDAO.ServiceLoadBalance{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	accessRepository := &serviceDAO.ServiceAccessControlRepository{}
	accessRule, err := accessRepository.Find(&serviceDAO.ServiceAccessControl{ServiceID: dt.ID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	return &serviceDAO.ServiceDetail{
		Info:          dt,
		Http:          httpRule,
		Tcp:           tcpRule,
		Grpc:          grpcRule,
		LoadBalance:   loadRule,
		AccessControl: accessRule,
	}, nil
}

func (s *ServiceInfoSvcLayer) PageList(info string, pageNum int32, pageSize int32) (int64, []serviceDAO.ServiceInfo, error) {
	var serviceList []serviceDAO.ServiceInfo

	query := dao.DB.Model(serviceDAO.ServiceInfo{})

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

func (s *ServiceInfoSvcLayer) NewHTTPService(req services.ServiceAddHTTPRequest) error {
	tx := dao.DB.Begin()

	serviceInfo := &serviceDAO.ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
	}

	result := tx.Create(&serviceInfo)
	if result.Error != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	httpRepository := &serviceDAO.ServiceHttpRuleRepository{}
	if err := httpRepository.Save(&serviceDAO.ServiceHttpRule{
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
		return status.Errorf(codes.InternalError, err.Error())
	}

	accessRepository := &serviceDAO.ServiceAccessControlRepository{}
	if err := accessRepository.Save(&serviceDAO.ServiceAccessControl{
		ServiceID:         serviceInfo.ID,
		OpenAuth:          int(req.OpenAuth),
		BlackList:         req.BlackList,
		WhiteList:         req.WhiteList,
		ClientIPFlowLimit: int(req.ClientIPFlowLimit),
		ServiceFlowLimit:  int(req.ServiceFlowLimit),
	}); err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	loadBalanceRepository := &serviceDAO.ServiceLoadBalanceRepository{}
	if err := loadBalanceRepository.Save(&serviceDAO.ServiceLoadBalance{
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
		return status.Errorf(codes.InternalError, err.Error())
	}

	tx.Commit()

	return nil
}
