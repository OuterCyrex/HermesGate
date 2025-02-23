package applicationDAO

import (
	"GoGateway/biz/model/application"
	"GoGateway/dao"
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model

	AppID    string `json:"app_id" gorm:"type:varchar(128);unique;not null"`
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Secret   string `json:"secret" gorm:"varchar(128);unique;not null"`
	WhiteIPS string `json:"white_ips" gorm:"varchar(800)"`
	Qpd      int64  `json:"qpd" gorm:"type:int"`
	Qps      int64  `json:"qps" gorm:"type:int"`
}

func (a *Application) TableName() string {
	return "go_gateway_application"
}

func (a *Application) ToHttpResponse() application.AppDetailResponse {
	return application.AppDetailResponse{
		ID:       int32(a.ID),
		AppID:    a.AppID,
		Name:     a.Name,
		Secret:   a.Secret,
		WhiteIPs: a.WhiteIPS,
		Qpd:      a.Qpd,
		QPS:      a.Qps,
	}
}

type ApplicationRepository struct{}

func (a *ApplicationRepository) Find(dt *Application) (*Application, error) {
	result := &Application{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
