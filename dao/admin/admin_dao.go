package adminDAO

import (
	"GoGateway/dao"
	"encoding/json"
	"gorm.io/gorm"
)

// Model

type Admin struct {
	gorm.Model

	Username string `json:"username"`
	Salt     string `json:"-"`
	Password string `json:"password"`
}

func (a *Admin) TableName() string {
	return "go_gateway_admin"
}

type AdminSessionInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	LoginTime int64  `json:"login_time"`
}

// DAO

type AdminRepository struct{}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{}
}

func (a *AdminRepository) Find(dt *Admin) (*Admin, error) {
	result := &Admin{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (asi *AdminSessionInfo) String() []byte {
	b, _ := json.Marshal(asi)

	return b
}

func (asi *AdminSessionInfo) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, asi)
	if err != nil {
		return err
	}
	return nil
}
