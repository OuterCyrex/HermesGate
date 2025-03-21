package serviceDAO

import (
	"GoGateway/dao"
)

type ServiceHttpRule struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	ServiceID      uint   `json:"service_id"`
	RuleType       int    `json:"rule_type" gorm:"type:int(11);default:0"`
	Rule           string `json:"rule" gorm:"type:varchar(300)"`
	NeedHttps      int    `json:"need_https" gorm:"type:tinyint"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"type:tinyint"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"type:tinyint"`
	UrlRewrite     string `json:"need_rewrite" gorm:"type:varchar(500)"`
	HeaderTransfer string `json:"header_transfer" gorm:"type:varchar(500)"`
}

func (s *ServiceHttpRule) TableName() string {
	return "go_gateway_service_http_rule"
}

// Repository Methods

type ServiceHttpRuleRepository struct{}

func (shr *ServiceHttpRuleRepository) Find(dt *ServiceHttpRule) (*ServiceHttpRule, error) {
	result := &ServiceHttpRule{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (shr *ServiceHttpRuleRepository) Save(dt *ServiceHttpRule) error {
	return dao.DB.Save(dt).Error
}
