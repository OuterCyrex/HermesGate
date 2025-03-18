package load_balance

import (
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/status"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashBalance struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int
	keys     UInt32Slice
	hashMap  map[uint32]string

	conf LoadBalanceConf
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		// crc32.ChecksumIEEE limits the node all in one hash space
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Add(param ...string) error {
	if len(param) == 0 {
		return status.Errorf(codes.InvalidParams, "at least 1 param is required")
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, addr := range param {
		// add nodes into hash space whose amount is the replicas
		for i := 0; i < c.replicas; i++ {
			hash := c.hash([]byte(strconv.Itoa(i) + addr))
			c.keys = append(c.keys, hash)
			c.hashMap[hash] = addr
		}
	}
	sort.Sort(c.keys)
	return nil
}

func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", status.Errorf(codes.NotFound, "no service node is in the hash space")
	}
	hash := c.hash([]byte(key))

	// search for the nearest node
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})
	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.hashMap[c.keys[idx]], nil
}

func (c *ConsistentHashBalance) SetConf(conf LoadBalanceConf) {
	c.conf = conf
}

func (c *ConsistentHashBalance) Update() {
	if conf, ok := c.conf.(*LoadBalanceCheckConf); ok {
		c.mux.Lock()
		defer c.mux.Unlock()
		c.keys = nil
		c.hashMap = map[uint32]string{}
		for _, ip := range conf.GetConf() {
			_ = c.Add(strings.Split(ip, ",")[0])
		}
	}
}
