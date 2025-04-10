// Code generated by hertz generator. DO NOT EDIT.

package services

import (
	services "GoGateway/biz/handler/services"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_service := root.Group("/service", _serviceMw()...)
		_service.GET("/list", append(_servicelistMw(), services.ServiceList)...)
		{
			_add := _service.Group("/add", _addMw()...)
			_add.POST("/grpc", append(_serviceaddgrpcMw(), services.ServiceAddGRPC)...)
			_add.POST("/http", append(_serviceaddhttpMw(), services.ServiceAddHTTP)...)
			_add.POST("/tcp", append(_serviceaddtcpMw(), services.ServiceAddTCP)...)
		}
		{
			_delete := _service.Group("/delete", _deleteMw()...)
			_delete.DELETE("/:id", append(_servicedeleteMw(), services.ServiceDelete)...)
		}
		{
			_detail := _service.Group("/detail", _detailMw()...)
			_detail.GET("/:id", append(_servicedetailMw(), services.ServiceDetail)...)
		}
		{
			_stat := _service.Group("/stat", _statMw()...)
			_stat.GET("/:id", append(_servicestaticMw(), services.ServiceStatic)...)
		}
		{
			_update := _service.Group("/update", _updateMw()...)
			{
				_grpc := _update.Group("/grpc", _grpcMw()...)
				_grpc.PUT("/:id", append(_serviceupdategrpcMw(), services.ServiceUpdateGRPC)...)
			}
			{
				_http := _update.Group("/http", _httpMw()...)
				_http.PUT("/:id", append(_serviceupdatehttpMw(), services.ServiceUpdateHTTP)...)
			}
			{
				_tcp := _update.Group("/tcp", _tcpMw()...)
				_tcp.PUT("/:id", append(_serviceupdatetcpMw(), services.ServiceUpdateTCP)...)
			}
		}
	}
}
