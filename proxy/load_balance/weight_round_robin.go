package load_balance

import (
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"strconv"
	"strings"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
	conf     LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int
	currentWeight   int
	effectiveWeight int
}

func (r *WeightRoundRobinBalance) Add(param ...string) error {
	if len(param) != 2 {
		return status.Errorf(codes.InvalidParams, "weight round robin should be 2 params")
	}
	parInt, err := strconv.ParseInt(param[1], 10, 64)
	if err != nil {
		return status.Errorf(codes.InternalError, err.Error())
	}
	node := &WeightNode{addr: param[0], weight: int(parInt), effectiveWeight: int(parInt)}
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		total += w.effectiveWeight

		w.currentWeight += w.effectiveWeight
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}

		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}

	best.currentWeight -= total
	return best.addr
}

func (r *WeightRoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *WeightRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *WeightRoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		r.rss = nil
		for _, ip := range conf.GetConf() {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}
}
