package applicationSVC

import (
	"GoGateway/biz/model/application"
	"GoGateway/dao"
	applicationDAO "GoGateway/dao/application"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"crypto/md5"
	"fmt"
	"io"
)

type ApplicationSvcLayer struct{}

func (a *ApplicationSvcLayer) NewApp(req application.AppAddHttpRequest) error {
	repo := applicationDAO.ApplicationRepository{}
	_, err := repo.Find(&applicationDAO.Application{AppID: req.AppID})
	if err == nil {
		return status.Errorf(codes.AlreadyExists, "AppID %s 已存在", req.AppID)
	}

	if req.Secret == "" {
		req.Secret = newSecretKey(req.AppID)
	}

	if err := dao.DB.Save(&applicationDAO.Application{
		AppID:    req.AppID,
		Name:     req.Name,
		Secret:   req.Secret,
		WhiteIPS: req.WhiteIPS,
		Qpd:      req.Qpd,
		Qps:      req.QPS,
	}).Error; err != nil {
		return status.Errorf(codes.InternalError, err.Error())
	}

	return nil
}

func newSecretKey(s string) string {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
