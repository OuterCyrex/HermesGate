// Code generated by hertz generator.

package services

import (
	services "GoGateway/biz/model/services"
	"GoGateway/conf"
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg/consts/codes"
	serviceConsts "GoGateway/pkg/consts/service"
	"GoGateway/pkg/status"
	serviceSVC "GoGateway/svc/services"
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// ServiceList .
// @router /service/list [GET]
func ServiceList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	svc := serviceSVC.ServiceInfoSvcLayer{}

	total, serviceList, err := svc.PageList(req.Info, req.PageNum, req.PageSize)
	if err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	var respList []*services.ServiceListItemResponse

	for _, s := range serviceList {
		serviceDetail, err := svc.ServiceDetail(&s)
		if err != nil {
			status.ErrToHttpResponse(c, err)
			return
		}

		serviceAddr := "unknown"

		clusterInfo := conf.GetConfig().Cluster

		switch serviceDetail.Info.LoadType {
		case serviceConsts.ServiceLoadTypeHTTP:
			port := clusterInfo.Port
			if serviceDetail.Http.NeedHttps == 1 {
				port = clusterInfo.SSLPort
			}
			if serviceDetail.Http.RuleType == serviceConsts.HTTPRuleTypePrefixURL {
				serviceAddr = fmt.Sprintf("%s:%d.%s", clusterInfo.IP, port, serviceDetail.Http.Rule)
			} else if serviceDetail.Http.RuleType == serviceConsts.HTTPRuleTypeDomain {
				serviceAddr = serviceDetail.Http.Rule
			}
		case serviceConsts.ServiceLoadTypeGRPC:
			serviceAddr = fmt.Sprintf("%s:%d", clusterInfo.IP, serviceDetail.Grpc.Port)
		case serviceConsts.ServiceLoadTypeTCP:
			serviceAddr = fmt.Sprintf("%s:%d", clusterInfo.IP, serviceDetail.Tcp.Port)
		default:
			hlog.Errorf("service type not support: %v", serviceDetail.Info.LoadType)
			c.JSON(http.StatusInternalServerError, "服务器内部错误")
		}

		ipList := strings.Split(serviceDetail.LoadBalance.IpList, ",")

		respList = append(respList, s.ToHttpResponse(serviceAddr, len(ipList)))
	}

	resp := services.ServiceListResponse{
		Total: total,
		Data:  respList,
	}

	c.JSON(consts.StatusOK, resp)
}

// ServiceDelete .
// @router /service/delete [DELETE]
func ServiceDelete(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceDeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	repository := serviceDAO.ServiceInfoRepository{}

	err = repository.Delete(uint(req.ID))
	if err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	resp := new(services.MessageResponse)

	resp.Message = "删除成功"

	c.JSON(consts.StatusOK, resp)
}

// ServiceAddHTTP .
// @router /service/add/http [POST]
func ServiceAddHTTP(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceAddHTTPRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	// validate input

	infoRepo := serviceDAO.ServiceInfoRepository{}
	if _, err = infoRepo.Find(&serviceDAO.ServiceInfo{ServiceName: req.ServiceName}); err == nil {
		c.JSON(http.StatusConflict, status.NewErrorResponse("服务名已被占用"))
		return
	}

	httpRepo := serviceDAO.ServiceHttpRuleRepository{}
	if _, err = httpRepo.Find(&serviceDAO.ServiceHttpRule{RuleType: int(req.RuleType), Rule: req.Rule}); err == nil {
		c.JSON(http.StatusConflict, status.NewErrorResponse("服务接入前缀或域名已存在"))
		return
	}

	if len(strings.Split(req.IpList, "\n")) != len(strings.Split(req.WeightList, "\n")) {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse("ip与权重列表不等"))
		return
	}

	svc := serviceSVC.ServiceInfoSvcLayer{}

	if err := svc.NewHTTPService(req); err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	resp := new(services.MessageResponse)

	resp.Message = "HTTP服务创建成功"

	c.JSON(consts.StatusOK, resp)
}

// ServiceUpdateHTTP .
// @router /service/update/:id [PUT]
func ServiceUpdateHTTP(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceUpdateHTTPRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	if len(strings.Split(req.IpList, "\n")) != len(strings.Split(req.WeightList, "\n")) {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse("ip与权重列表不等"))
		return
	}

	svc := serviceSVC.ServiceInfoSvcLayer{}
	if err := svc.UpdateHTTPService(req); err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	resp := new(services.MessageResponse)

	resp.Message = "HTTP服务更新成功"

	c.JSON(consts.StatusOK, resp)
}

// ServiceDetail .
// @router /service/detail/:id [GET]
func ServiceDetail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceDetailRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	svc := serviceSVC.ServiceInfoSvcLayer{}

	repo := serviceDAO.ServiceInfoRepository{}

	info, err := repo.Find(&serviceDAO.ServiceInfo{Model: gorm.Model{ID: uint(req.ID)}})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		status.ErrToHttpResponse(c, status.Errorf(codes.NotFound, "服务信息不存在"))
		return
	} else if err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	detail, err := svc.ServiceDetail(info)
	if err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	c.JSON(consts.StatusOK, detail.ToHttpResponse())
}

// ServiceStatic .
// @router /service/static/:id [GET]
func ServiceStatic(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceStaticRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	var today []int64
	var yesterday []int64

	for i := 0; i <= time.Now().Hour(); i++ {
		today = append(today, 0)
	}

	for i := 0; i <= 23; i++ {
		yesterday = append(yesterday, 0)
	}

	c.JSON(consts.StatusOK, services.ServiceStaticResponse{
		Today:     today,
		Yesterday: yesterday,
	})
}

// ServiceAddGRPC .
// @router /service/add/grpc [POST]
func ServiceAddGRPC(ctx context.Context, c *app.RequestContext) {
	var err error
	var req services.ServiceAddGrpcRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, status.NewErrorResponse(err.Error()))
		return
	}

	svc := serviceSVC.ServiceInfoSvcLayer{}
	if err := svc.NewGRPCService(req); err != nil {
		status.ErrToHttpResponse(c, err)
		return
	}

	resp := new(services.MessageResponse)

	resp.Message = "GRPC服务创建成功"

	c.JSON(consts.StatusOK, resp)
}
