package main

import (
	"GoGateway/dao"
	adminDAO "GoGateway/dao/admin"
	serviceDAO "GoGateway/dao/services"
)

func main() {
	dao.InitDB(dao.DefaultDSN())

	err := dao.DB.AutoMigrate(
		&adminDAO.Admin{},
		&serviceDAO.ServiceInfo{},
		&serviceDAO.ServiceHttpRule{},
		&serviceDAO.ServiceGRPCRule{},
		&serviceDAO.ServiceTcpRule{},
		&serviceDAO.ServiceLoadBalance{},
		&serviceDAO.ServiceAccessControl{},
	)
	if err != nil {
		panic(err)
	}
}
