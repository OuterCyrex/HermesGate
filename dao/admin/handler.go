package adminDAO

import (
	"GoGateway/dao"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"gorm.io/gorm"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (a *AdminHandler) Find(dt *Admin) (*Admin, error) {
	result := &Admin{}
	err := dao.DB.Table(dt.TableName()).Where(dt).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AdminHandler) LoginAndCheck(param *Admin) (*Admin, error) {
	adminInfo, err := a.Find(&Admin{Username: param.Username})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户信息不存在")
	}
	saltPassword := genSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, status.Errorf(codes.InvalidParams, "密码错误")
	}

	return adminInfo, nil
}

func (a *AdminHandler) UpdatePassword(Id uint, password string) error {
	adminInfo, err := a.Find(&Admin{Model: gorm.Model{ID: Id}})
	if err != nil {
		return status.Errorf(codes.NotFound, "用户信息不存在")
	}
	saltPassword := genSaltPassword(adminInfo.Salt, password)

	result := dao.DB.Table(adminInfo.TableName()).Where(&Admin{Model: gorm.Model{ID: Id}}).Updates(&Admin{Password: saltPassword})
	if result.Error != nil {
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	return nil
}
