package load_balance

type LoadBalanceType int

const (
	LoadBalanceRandom LoadBalanceType = iota
	LoadBalanceRoundRobin
	LoadBalanceWeightRoundRobin
	LoadBalanceConsistentHash
)

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	Update()
}

func LoadBalanceFactory(t LoadBalanceType) LoadBalance {
	switch t {
	case LoadBalanceRandom:
		return &RandomBalance{}
	case LoadBalanceRoundRobin:
		return &RoundRobinBalance{}
	case LoadBalanceWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	case LoadBalanceConsistentHash:
		return &ConsistentHashBalance{}
	}
	return &RoundRobinBalance{}
}

func LoadBalanceFactorWithConf(t LoadBalanceType, lbConf LoadBalanceConf) LoadBalance {
	//观察者模式
	switch t {
	case LoadBalanceRandom:
		lb := &RandomBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		lb.Update()
		return lb
	case LoadBalanceConsistentHash:
		lb := NewConsistentHashBalance(10, nil)
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		lb.Update()
		return lb
	case LoadBalanceRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		lb.Update()
		return lb
	case LoadBalanceWeightRoundRobin:
		lb := &WeightRoundRobinBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		lb.Update()
		return lb
	default:
		lb := &RoundRobinBalance{}
		lb.SetConf(lbConf)
		lbConf.Attach(lb)
		lb.Update()
		return lb
	}
}
