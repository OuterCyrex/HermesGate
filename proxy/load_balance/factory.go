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
	case LoadBalanceRoundRobin:
	case LoadBalanceWeightRoundRobin:
	case LoadBalanceConsistentHash:
	}
	return nil
}
