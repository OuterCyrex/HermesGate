package dashboardSVC

import (
	"GoGateway/biz/model/dashboard"
	"GoGateway/dao"
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/consts/codes"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/pkg/status"
)

type DashboardSvcLayer struct{}

func (d *DashboardSvcLayer) GetServiceStat() (*dashboard.DashServiceStatResponse, error) {
	type Output struct {
		LoadType int64
		Count    int64
	}

	var results []Output
	var respList []*dashboard.DashServiceStatItem

	err := dao.DB.Model(&serviceDAO.ServiceInfo{}).
		Select("load_type, COUNT(*) as count").
		Group("load_type").Find(&results).Error
	if err != nil {
		return nil, status.Errorf(codes.InternalError, err.Error())
	}

	for _, r := range results {
		name := "unknown"

		switch r.LoadType {
		case serviceConsts.ServiceLoadTypeHTTP:
			name = "http"
		case serviceConsts.ServiceLoadTypeTCP:
			name = "tcp"
		case serviceConsts.ServiceLoadTypeGRPC:
			name = "grpc"
		}

		respList = append(respList, &dashboard.DashServiceStatItem{
			Name:  name,
			Value: r.Count,
		})
	}

	return &dashboard.DashServiceStatResponse{
		Total: int64(len(results)),
		Data:  respList,
	}, nil
}
