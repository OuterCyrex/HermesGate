package adminDAO

import (
	"encoding/json"
	"gorm.io/gorm"
)

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
