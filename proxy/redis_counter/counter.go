package redisCounter

import (
	"GoGateway/pkg"
	"GoGateway/pkg/consts/codes"
	"GoGateway/pkg/consts/redis"
	"GoGateway/pkg/status"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type ServiceFlowCounter struct {
	counterMap map[string]*RedisCounter
	Mutex      sync.Mutex
}

func newServiceFlowCounter() *ServiceFlowCounter {
	return &ServiceFlowCounter{
		counterMap: make(map[string]*RedisCounter),
		Mutex:      sync.Mutex{},
	}
}

func (h *ServiceFlowCounter) GetAllInfo() (int64, int64) {
	var total int64
	var totalQPS int64
	for _, v := range h.counterMap {
		total += v.TotalCount
		totalQPS += v.QPS
	}
	return total, totalQPS
}

var ServiceFlowCountHandler *ServiceFlowCounter

func init() {
	ServiceFlowCountHandler = newServiceFlowCounter()
}

func (h *ServiceFlowCounter) GetCounter(serviceName string) *RedisCounter {
	if v, ok := h.counterMap[serviceName]; ok {
		return v
	}

	counter := NewRedisCounter(serviceName, 1*time.Second)
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.counterMap[serviceName] = counter

	return counter
}

type RedisCounter struct {
	ServiceName     string
	Interval        time.Duration
	QPS             int64
	TickerBeginTime int64
	TickerCount     int64
	TotalCount      int64
}

func NewRedisCounter(serviceName string, interval time.Duration) *RedisCounter {
	counter := &RedisCounter{
		ServiceName:     serviceName,
		Interval:        interval,
		QPS:             0,
		TickerBeginTime: time.Now().Unix(),
	}

	// 启动监控线程
	go func() {
		// 初始化计时器
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&counter.TickerCount)
			atomic.StoreInt64(&counter.TickerCount, 0)
			formerTotalCount := atomic.LoadInt64(&counter.TotalCount)
			lastTickerBeginTime := atomic.LoadInt64(&counter.TickerBeginTime)

			// 获取redis数据
			current := time.Now()
			dayKey := counter.DayKey(current)
			hourKey := counter.HourKey(current)
			totalHourKey := counter.TotalHourKey(current)

			// 开启 pipeline
			_, err := Pipeline(
				counter.send("INCRBY", dayKey, tickerCount),
				counter.send("EXPIRE", dayKey, 60*60*24*2),
				counter.send("INCRBY", hourKey, tickerCount),
				counter.send("EXPIRE", hourKey, 60*60*24*2),
				counter.send("INCRBY", totalHourKey, tickerCount),
				counter.send("EXPIRE", totalHourKey, 60*60*24*2),
			)
			if err != nil {
				hlog.Errorf("redis pipeline failed, %s", err.Error())
				continue
			}

			totalCount, err := counter.DayCount(current)
			nowUnix := time.Now().Unix()
			tickerCount = totalCount - formerTotalCount
			if nowUnix > lastTickerBeginTime {
				atomic.StoreInt64(&counter.QPS, tickerCount/(nowUnix-lastTickerBeginTime))
				atomic.StoreInt64(&counter.TotalCount, formerTotalCount+tickerCount)
				atomic.StoreInt64(&counter.TickerBeginTime, nowUnix)
			}
		}
	}()

	return counter
}

func (c *RedisCounter) DayCount(current time.Time) (int64, error) {
	conn := pkg.GetRedis()
	defer func() {
		_ = conn.Close()
	}()
	result, err := conn.Do("GET", c.DayKey(current))
	v, ok := result.([]byte)
	if !ok {
		return 0, status.Errorf(codes.InternalError, "Redis Internal Error")
	}
	count, _ := strconv.Atoi(string(v))
	return int64(count), err
}
func (c *RedisCounter) HourCount(current time.Time) (int64, error) {
	conn := pkg.GetRedis()
	defer func() {
		_ = conn.Close()
	}()
	result, err := conn.Do("GET", c.HourKey(current))
	v, ok := result.([]byte)
	if !ok {
		return 0, status.Errorf(codes.InternalError, "Redis Internal Error")
	}
	count, _ := strconv.Atoi(string(v))
	return int64(count), err
}

func (c *RedisCounter) DayKey(current time.Time) string {
	return fmt.Sprintf("%s_%s_%s", redisConsts.RedisCounterKey, c.ServiceName, current.Format("20060102"))
}
func (c *RedisCounter) HourKey(current time.Time) string {
	return fmt.Sprintf("%s_%s_%s", redisConsts.RedisCounterKey, c.ServiceName, current.Format("2006010215"))
}

func (c *RedisCounter) TotalHourKey(current time.Time) string {
	return fmt.Sprintf("%s_%s_%s", redisConsts.RedisCounterTotalKey, current.Format("2006010215"))
}

func (c *RedisCounter) TotalHourCount(current time.Time) (int64, error) {
	conn := pkg.GetRedis()
	defer func() {
		_ = conn.Close()
	}()
	result, err := conn.Do("GET", c.TotalHourKey(current))
	v, ok := result.([]byte)
	if !ok {
		return 0, status.Errorf(codes.InternalError, "Redis Internal Error")
	}
	count, _ := strconv.Atoi(string(v))
	return int64(count), err
}

// Increase 异步原子增加
func (c *RedisCounter) Increase() {
	go func() {
		atomic.AddInt64(&c.TickerCount, 1)
	}()
}

// 对 Pipeline 细节进行封装
func (c *RedisCounter) send(commandName string, args ...interface{}) sendInput {
	return sendInput{
		commandName: commandName,
		args:        args,
	}
}

type sendInput struct {
	commandName string
	args        []interface{}
}

func Pipeline(sends ...sendInput) ([]interface{}, error) {
	c := pkg.GetRedis()

	for _, send := range sends {
		err := c.Send(send.commandName, send.args...)
		if err != nil {
			_ = c.Close()
			return nil, err
		}
	}

	err := c.Flush()
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	err = c.Close()
	return nil, err
}
