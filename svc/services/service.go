package serviceSVC

import (
	"GoGateway/biz/model/services"
	"GoGateway/dao"
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/consts/codes"
	serviceConsts "GoGateway/pkg/consts/service"
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
		LoadType:    serviceConsts.ServiceLoadTypeHTTP,
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
		NeedHttps:      int(req.NeedHTTPS),
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

func (s *ServiceInfoSvcLayer) NewTCPService(req services.ServiceAddTcpRequest) error {
	tx := dao.DB.Begin()

	serviceInfo := &serviceDAO.ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
		LoadType:    serviceConsts.ServiceLoadTypeTCP,
	}

	result := tx.Create(&serviceInfo)
	if result.Error != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	httpRepository := &serviceDAO.ServiceTcpRuleRepository{}
	if err := httpRepository.Save(&serviceDAO.ServiceTcpRule{
		ServiceID: serviceInfo.ID,
		Port:      int(req.Port),
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
		ServiceID:  serviceInfo.ID,
		RoundType:  int(req.RoundType),
		IpList:     req.IpList,
		WeightList: req.WeightList,
		ForbidList: req.ForbidList,
	}); err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	tx.Commit()

	return nil
}

func (s *ServiceInfoSvcLayer) NewGRPCService(req services.ServiceAddGrpcRequest) error {
	tx := dao.DB.Begin()

	serviceInfo := &serviceDAO.ServiceInfo{
		ServiceName: req.ServiceName,
		ServiceDesc: req.ServiceDesc,
		LoadType:    serviceConsts.ServiceLoadTypeGRPC,
	}

	result := tx.Create(&serviceInfo)
	if result.Error != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	httpRepository := &serviceDAO.ServiceGRPCRuleRepository{}
	if err := httpRepository.Save(&serviceDAO.ServiceGRPCRule{
		ServiceID:      serviceInfo.ID,
		Port:           int(req.Port),
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
		ServiceID:  serviceInfo.ID,
		RoundType:  int(req.RoundType),
		IpList:     req.IpList,
		WeightList: req.WeightList,
		ForbidList: req.ForbidList,
	}); err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	tx.Commit()

	return nil
}

func (s *ServiceInfoSvcLayer) UpdateHTTPService(req services.ServiceUpdateHTTPRequest) error {
	tx := dao.DB.Begin()
	serviceDetail, err := s.ServiceDetail(&serviceDAO.ServiceInfo{Model: gorm.Model{ID: uint(req.ID)}})
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// httpRule
	if err := tx.Save(&serviceDAO.ServiceHttpRule{
		ID:             serviceDetail.Http.ID,
		ServiceID:      serviceDetail.Http.ServiceID,
		RuleType:       serviceDetail.Http.RuleType,
		Rule:           serviceDetail.Http.Rule,
		NeedHttps:      int(req.NeedHTTPS),
		NeedWebsocket:  int(req.NeedWebsocket),
		NeedStripUri:   int(req.NeedStripUri),
		UrlRewrite:     req.UrlRewrite,
		HeaderTransfer: req.HeaderTransfer,
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// access control
	if err := tx.Save(&serviceDAO.ServiceAccessControl{
		ID:                serviceDetail.AccessControl.ID,
		ServiceID:         serviceDetail.AccessControl.ServiceID,
		OpenAuth:          int(req.OpenAuth),
		BlackList:         req.BlackList,
		WhiteList:         req.WhiteList,
		WhiteHostName:     serviceDetail.AccessControl.WhiteHostName,
		ClientIPFlowLimit: int(req.ClientIPFlowLimit),
		ServiceFlowLimit:  int(req.ServiceFlowLimit),
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// load balance
	if err := tx.Save(&serviceDAO.ServiceLoadBalance{
		ID:                     serviceDetail.LoadBalance.ID,
		ServiceID:              serviceDetail.LoadBalance.ServiceID,
		CheckMethod:            serviceDetail.LoadBalance.CheckMethod,
		CheckTimeout:           serviceDetail.LoadBalance.CheckTimeout,
		CheckInterval:          serviceDetail.LoadBalance.CheckInterval,
		RoundType:              int(req.RoundType),
		IpList:                 req.IpList,
		WeightList:             req.WeightList,
		ForbidList:             serviceDetail.LoadBalance.ForbidList,
		UpstreamConnectTimeout: int(req.UpstreamConnectTimeout),
		UpstreamHeaderTimeout:  int(req.UpstreamHeaderTimeout),
		UpstreamIdleTimeout:    int(req.UpstreamIdleTimeout),
		UpstreamMaxIdle:        int(req.UpstreamMaxIdle),
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}
	return nil
}

func (s *ServiceInfoSvcLayer) UpdateGrpcService(req services.ServiceUpdateGrpcRequest) error {
	tx := dao.DB.Begin()
	serviceDetail, err := s.ServiceDetail(&serviceDAO.ServiceInfo{Model: gorm.Model{ID: uint(req.ID)}})
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// httpRule
	if err := tx.Save(&serviceDAO.ServiceGRPCRule{
		ID:             serviceDetail.Grpc.ID,
		ServiceID:      serviceDetail.Grpc.ServiceID,
		Port:           serviceDetail.Grpc.Port,
		HeaderTransfer: req.HeaderTransfer,
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// access control
	if err := tx.Save(&serviceDAO.ServiceAccessControl{
		ID:                serviceDetail.AccessControl.ID,
		ServiceID:         serviceDetail.AccessControl.ServiceID,
		OpenAuth:          int(req.OpenAuth),
		BlackList:         req.BlackList,
		WhiteList:         req.WhiteList,
		WhiteHostName:     req.WhiteHostName,
		ClientIPFlowLimit: int(req.ClientIPFlowLimit),
		ServiceFlowLimit:  int(req.ServiceFlowLimit),
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// load balance
	if err := tx.Save(&serviceDAO.ServiceLoadBalance{
		ID:            serviceDetail.LoadBalance.ID,
		ServiceID:     serviceDetail.LoadBalance.ServiceID,
		CheckMethod:   serviceDetail.LoadBalance.CheckMethod,
		CheckTimeout:  serviceDetail.LoadBalance.CheckTimeout,
		CheckInterval: serviceDetail.LoadBalance.CheckInterval,
		RoundType:     int(req.RoundType),
		IpList:        req.IpList,
		WeightList:    req.WeightList,
		ForbidList:    req.ForbidList,
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}
	return nil
}

func (s *ServiceInfoSvcLayer) UpdateTcpService(req services.ServiceUpdateTcpRequest) error {
	tx := dao.DB.Begin()
	serviceDetail, err := s.ServiceDetail(&serviceDAO.ServiceInfo{Model: gorm.Model{ID: uint(req.ID)}})
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// tcpRule
	if err := tx.Save(&serviceDAO.ServiceTcpRule{
		ID:        serviceDetail.Tcp.ID,
		ServiceID: serviceDetail.Tcp.ServiceID,
		Port:      serviceDetail.Tcp.Port,
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// access control
	if err := tx.Save(&serviceDAO.ServiceAccessControl{
		ID:                serviceDetail.AccessControl.ID,
		ServiceID:         serviceDetail.AccessControl.ServiceID,
		OpenAuth:          int(req.OpenAuth),
		BlackList:         req.BlackList,
		WhiteList:         req.WhiteList,
		WhiteHostName:     req.WhiteHostName,
		ClientIPFlowLimit: int(req.ClientIPFlowLimit),
		ServiceFlowLimit:  int(req.ServiceFlowLimit),
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	// load balance
	if err := tx.Save(&serviceDAO.ServiceLoadBalance{
		ID:            serviceDetail.LoadBalance.ID,
		ServiceID:     serviceDetail.LoadBalance.ServiceID,
		CheckMethod:   serviceDetail.LoadBalance.CheckMethod,
		CheckTimeout:  serviceDetail.LoadBalance.CheckTimeout,
		CheckInterval: serviceDetail.LoadBalance.CheckInterval,
		RoundType:     int(req.RoundType),
		IpList:        req.IpList,
		WeightList:    req.WeightList,
		ForbidList:    req.ForbidList,
	}).Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return status.Errorf(codes.InternalError, err.Error())
	}
	return nil
}
