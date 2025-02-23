package applicationSVC

import (
	"GoGateway/biz/model/application"
	"GoGateway/dao"
	applicationDAO "GoGateway/dao/application"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"io"
	"time"
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

func (a *ApplicationSvcLayer) UpdateApp(req application.AppUpdateRequest) error {
	repo := applicationDAO.ApplicationRepository{}
	info, err := repo.Find(&applicationDAO.Application{
		Model: gorm.Model{ID: uint(req.ID)},
	})
	if err != nil {
		return status.Errorf(codes.NotFound, "App信息不存在")
	}

	secret := newSecretKey(info.AppID)

	if req.Secret != "" {
		secret = req.Secret
	}

	if err := repo.Save(&applicationDAO.Application{
		Model: gorm.Model{
			ID:        uint(req.ID),
			CreatedAt: info.CreatedAt,
			UpdatedAt: time.Now(),
		},
		AppID:    info.AppID,
		Name:     req.Name,
		Secret:   secret,
		WhiteIPS: req.WhiteIPS,
		Qpd:      req.Qpd,
		Qps:      req.QPS,
	}); err != nil {
		return status.Errorf(codes.InternalError, err.Error())
	}

	return nil
}

func (a *ApplicationSvcLayer) DeleteApp(req application.AppDeleteRequest) error {
	result := dao.DB.Delete(&applicationDAO.Application{Model: gorm.Model{ID: uint(req.ID)}})
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "App信息不存在")
	}
	if result.Error != nil {
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	return nil
}

func (a *ApplicationSvcLayer) AppList(req application.AppListRequest) (*application.AppListResponse, error) {
	var apps []applicationDAO.Application
	query := dao.DB.Model(&applicationDAO.Application{})

	if req.Info != "" {
		query = query.Where("name like ? or app_id like ?", "%"+req.Info+"%", "%"+req.Info+"%")
	}

	var count int64

	query.Count(&count)

	result := query.Scopes(dao.Paginate(req.PageNum, req.PageSize)).Find(&apps)
	if result.Error != nil {
		return nil, status.Errorf(codes.InternalError, result.Error.Error())
	}

	var respList []*application.AppListItemResponse
	for _, app := range apps {
		respList = append(respList, &application.AppListItemResponse{
			ID:       int32(app.ID),
			AppID:    app.AppID,
			Name:     app.Name,
			Secret:   app.Secret,
			WhiteIPs: app.WhiteIPS,
			Qpd:      app.Qpd,
			QPS:      app.Qps,
			RealQps:  0,
			RealQpd:  0,
		})
	}

	return &application.AppListResponse{
		Total: count,
		Data:  respList,
	}, nil
}

func newSecretKey(s string) string {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
