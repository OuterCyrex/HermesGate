package load_balance

import (
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"strings"
)

type RoundRobinBalance struct {
	curIndex int
	rss      []string
	conf     LoadBalanceConf
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return status.Errorf(codes.InvalidParams, "at least 1 param is required")
	}
	r.rss = append(r.rss, params...)
	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	length := len(r.rss)
	if r.curIndex >= length {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % length
	return curAddr
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

// Update Inserts the loadBalanceConf into balance instance
func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		r.rss = nil
		for _, ip := range conf.GetConf() {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}
}
