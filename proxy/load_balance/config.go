package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"sync"
	"time"
)

type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type Observer interface {
	Update()
}

const (
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

// LoadBalanceCheckConf is the implement of LoadBalanceConf
type LoadBalanceCheckConf struct {
	observers []Observer

	// confIPWeight is the ip-weight key-value map
	confIPWeight map[string]string

	// activeList is a list of services that are active now
	activeList []string
	format     string
}

// NotifyAllObservers Updates all observers
func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, o := range s.observers {
		o.Update()
	}
}

func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceCheckConf) GetConf() []string {
	var confList []string
	for _, ip := range s.activeList {
		weight, ok := s.confIPWeight[ip]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, o := range s.observers {
		o.Update()
	}
}

func NewLoadBalanceCheckConf(format string, conf map[string]string) *LoadBalanceCheckConf {
	var activeList []string
	for item := range conf {
		activeList = append(activeList, item)
	}
	lbConf := &LoadBalanceCheckConf{
		format:       format,
		confIPWeight: conf,
		activeList:   activeList,
	}
	return lbConf
}

func (s *LoadBalanceCheckConf) WatchConf() {
	go func() {
		// confIPErrNum records the error nums of each
		confIPErrNum := make(map[string]int)
		for {
			var changedList []string
			var wg sync.WaitGroup
			for item := range s.confIPWeight {
				wg.Add(1)
				go func(ip string) {
					defer wg.Done()
					conn, err := net.DialTimeout("tcp", ip, DefaultCheckTimeout)
					if err == nil {
						_ = conn.Close()
						confIPErrNum[ip] = 0
					} else {
						confIPErrNum[ip] += 1
					}

					if confIPErrNum[ip] < DefaultCheckMaxErrNum {
						changedList = append(changedList, ip)
					}
				}(item)

				// if the failed num has been greater than MaxErrNum
				// remove the item from the activeList
				if confIPErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}

			// wait for all goroutines has finished the task
			wg.Wait()

			sort.Strings(changedList)
			sort.Strings(s.activeList)
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}

			// wait for the next turn of health check
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}
