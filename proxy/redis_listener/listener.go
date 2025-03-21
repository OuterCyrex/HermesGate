package reloadListener

import (
	"GoGateway/dao"
	serviceDAO "GoGateway/dao/services"
	"GoGateway/pkg"
	redisConsts "GoGateway/pkg/consts/redis"
	"GoGateway/proxy"
	serviceSVC "GoGateway/svc/services"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type reloader interface {
	Reload(*serviceDAO.ServiceDetail)
}
type ReloadListener struct {
	reloaderMap map[int]reloader
	lock        sync.Mutex
	balancer    *proxy.ServiceBalancer
	conn        *redis.PubSubConn
}

func NewReloadListener(
	balancer *proxy.ServiceBalancer,
) *ReloadListener {
	conn := pkg.GetRedis()
	psConn := redis.PubSubConn{Conn: conn}

	err := psConn.Subscribe(redisConsts.RedisChannelKey)
	if err != nil {
		panic(err)
		return nil
	}

	return &ReloadListener{
		reloaderMap: make(map[int]reloader),
		lock:        sync.Mutex{},
		balancer:    balancer,
		conn:        &psConn,
	}
}

func (l *ReloadListener) Add(loadType int, r reloader) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.reloaderMap[loadType] = r
}

func (l *ReloadListener) Listen() {
	for {
		switch v := l.conn.Receive().(type) {
		case redis.Message:
			id, err := strconv.Atoi(string(v.Data))
			if err != nil {
				hlog.Error("ReloadListener error: ", v)
				continue
			}
			detail, err := l.Get(uint(id))
			l.reloaderMap[detail.Info.LoadType].Reload(detail)

			err = l.balancer.ReloadLoadBalance(detail)
			if err != nil {
				hlog.Error("ReloadListener error: ", v)
				continue
			}
		case redis.Subscription:
			hlog.Infof("Subscription received: %s %s", v.Channel, v.Kind)
		case error:
			hlog.Error("ReloadListener error: ", v)
		}
	}
}

func (l *ReloadListener) Get(id uint) (*serviceDAO.ServiceDetail, error) {
	var info serviceDAO.ServiceInfo
	dao.DB.Model(&serviceDAO.ServiceInfo{}).Where(&serviceDAO.ServiceInfo{Model: gorm.Model{ID: id}}).First(&info)
	svc := serviceSVC.ServiceInfoSvcLayer{}
	return svc.ServiceDetail(&info)
}
