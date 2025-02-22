package serviceDAO

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info"`
	Http          *ServiceHttpRule      `json:"http"`
	Tcp           *ServiceTcpRule       `json:"tcp"`
	Grpc          *ServiceGRPCRule      `json:"grpc"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance"`
	AccessControl *ServiceAccessControl `json:"access_control"`
}
