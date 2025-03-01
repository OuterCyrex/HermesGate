package load_balance

import (
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"math/rand"
	"strings"
)

type RandomBalance struct {
	curIndex int
	rss      []string
	conf     LoadBalanceConf
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return status.Errorf(codes.InvalidParams, "at least 1 param is required")
	}
	r.rss = append(r.rss, params...)
	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *RandomBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		r.rss = nil
		for _, ip := range conf.GetConf() {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}
}
