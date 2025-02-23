package adminSVC

import (
	"GoGateway/dao"
	adminDAO "GoGateway/dao/admin"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"crypto/sha256"
	"fmt"
	"gorm.io/gorm"
)

type AdminSvcLayer struct{}

func (a *AdminSvcLayer) LoginAndCheck(param *adminDAO.Admin) (*adminDAO.Admin, error) {
	repository := adminDAO.NewAdminRepository()

	adminInfo, err := repository.Find(&adminDAO.Admin{Username: param.Username})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户信息不存在")
	}
	saltPassword := genSaltPassword(adminInfo.Salt, param.Password)

	if adminInfo.Password != saltPassword {
		return nil, status.Errorf(codes.InvalidParams, "密码错误")
	}

	return adminInfo, nil
}

func (a *AdminSvcLayer) UpdatePassword(Id uint, password string) error {
	repository := adminDAO.NewAdminRepository()

	adminInfo, err := repository.Find(&adminDAO.Admin{Model: gorm.Model{ID: Id}})
	if err != nil {
		return status.Errorf(codes.NotFound, "用户信息不存在")
	}
	saltPassword := genSaltPassword(adminInfo.Salt, password)

	result := dao.DB.Table(adminInfo.TableName()).Where(&adminDAO.Admin{Model: gorm.Model{ID: Id}}).Updates(&adminDAO.Admin{Password: saltPassword})
	if result.Error != nil {
		return status.Errorf(codes.InternalError, result.Error.Error())
	}

	return nil
}

func genSaltPassword(salt string, password string) string {
	sh1 := sha256.New()
	sh1.Write([]byte(salt))
	sh2 := sha256.New()
	sh2.Write([]byte(fmt.Sprintf("%x", sh1.Sum(nil)) + password))
	return fmt.Sprintf("%x", sh2.Sum(nil))
}
