package serviceDAO

import "GoGateway/biz/model/services"

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info"`
	Http          *ServiceHttpRule      `json:"http"`
	Tcp           *ServiceTcpRule       `json:"tcp"`
	Grpc          *ServiceGRPCRule      `json:"grpc"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance"`
	AccessControl *ServiceAccessControl `json:"access_control"`
}

func (sd *ServiceDetail) ToHttpResponse() *services.ServiceDetailResponse {

	detail := &services.ServiceDetailResponse{
		Info:          nil,
		Http:          nil,
		Tcp:           nil,
		Grpc:          nil,
		LoadBalance:   nil,
		AccessControl: nil,
	}

	if sd.Info != nil {
		detail.Info = &services.ServiceInfoPart{
			ID:          int32(sd.Info.ID),
			LoadType:    int8(sd.Info.LoadType),
			ServiceName: sd.Info.ServiceName,
			ServiceDesc: sd.Info.ServiceDesc,
		}
	}

	if sd.Http != nil {
		detail.Http = &services.ServiceHttpRulePart{
			ID:             int32(sd.Http.ID),
			ServiceID:      int32(sd.Http.ServiceID),
			RuleType:       int8(sd.Http.RuleType),
			Rule:           sd.Http.Rule,
			NeedHttps:      int8(sd.Http.NeedHttps),
			NeedWebsocket:  int8(sd.Http.NeedWebsocket),
			NeedStripUri:   int8(sd.Http.NeedStripUri),
			UrlRewrite:     sd.Http.UrlRewrite,
			HeaderTransfer: sd.Http.HeaderTransfer,
		}
	}

	if sd.Tcp != nil {
		detail.Tcp = &services.ServiceTcpRulePart{
			ID:        int32(sd.Tcp.ID),
			ServiceID: int32(sd.Tcp.ServiceID),
			Port:      int32(sd.Tcp.Port),
		}
	}

	if sd.Grpc != nil {
		detail.Grpc = &services.ServiceGRPCRulePart{
			ID:          int32(sd.Grpc.ID),
			ServiceID:   int32(sd.Grpc.ServiceID),
			Port:        int32(sd.Grpc.Port),
			HeaderTrans: sd.Grpc.HeaderTransfer,
		}
	}

	if sd.LoadBalance != nil {
		detail.LoadBalance = &services.ServiceLoadBalancePart{
			ID:                     int32(sd.LoadBalance.ID),
			ServiceID:              int32(sd.LoadBalance.ServiceID),
			CheckMethod:            int32(sd.LoadBalance.CheckMethod),
			CheckTimeout:           int32(sd.LoadBalance.CheckTimeout),
			CheckInterval:          int32(sd.LoadBalance.CheckInterval),
			RoundType:              int8(sd.LoadBalance.RoundType),
			IpList:                 sd.LoadBalance.IpList,
			WeightList:             sd.LoadBalance.WeightList,
			ForbidList:             sd.LoadBalance.ForbidList,
			UpstreamConnectTimeout: int32(sd.LoadBalance.UpstreamConnectTimeout),
			UpstreamHeaderTimeout:  int32(sd.LoadBalance.UpstreamHeaderTimeout),
			UpstreamIdleTimeout:    int32(sd.LoadBalance.UpstreamIdleTimeout),
			UpstreamMaxIdle:        int32(sd.LoadBalance.UpstreamMaxIdle),
		}
	}

	if sd.AccessControl != nil {
		detail.AccessControl = &services.ServiceAccessControlPart{
			ID:                int32(sd.AccessControl.ID),
			ServiceID:         int32(sd.AccessControl.ServiceID),
			OpenAuth:          int8(sd.AccessControl.OpenAuth),
			BlackList:         sd.AccessControl.BlackList,
			WhiteList:         sd.AccessControl.WhiteList,
			WhiteHostName:     sd.AccessControl.WhiteHostName,
			ClientIPFlowLimit: int32(sd.AccessControl.ClientIPFlowLimit),
			ServiceFlowLimit:  int32(sd.AccessControl.ServiceFlowLimit),
		}
	}

	return detail
}
