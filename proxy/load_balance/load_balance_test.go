package load_balance

import (
	"fmt"
	"testing"
)

func TestRandomBalance(t *testing.T) {
	rb := &RandomBalance{}
	_ = rb.Add("127.0.0.1:2003") //0
	_ = rb.Add("127.0.0.1:2004") //1
	_ = rb.Add("127.0.0.1:2005") //2
	_ = rb.Add("127.0.0.1:2006") //3
	_ = rb.Add("127.0.0.1:2007") //4

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
}

func TestRoundRobin(t *testing.T) {
	rb := &RoundRobinBalance{}
	_ = rb.Add("127.0.0.1:2003") //0
	_ = rb.Add("127.0.0.1:2004") //1
	_ = rb.Add("127.0.0.1:2005") //2
	_ = rb.Add("127.0.0.1:2006") //3
	_ = rb.Add("127.0.0.1:2007") //4

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
}

func TestConsistentHashBalance(t *testing.T) {
	rb := NewConsistentHashBalance(10, nil)
	_ = rb.Add("127.0.0.1:2003") //0
	_ = rb.Add("127.0.0.1:2004") //1
	_ = rb.Add("127.0.0.1:2005") //2
	_ = rb.Add("127.0.0.1:2006") //3
	_ = rb.Add("127.0.0.1:2007") //4

	//url hash
	fmt.Println(rb.Get("http://127.0.0.1:2002/user/1"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/user/2"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/user/3"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/user/4"))

	//ip hash
	fmt.Println(rb.Get("127.0.0.1"))
	fmt.Println(rb.Get("192.168.0.1"))
	fmt.Println(rb.Get("127.0.0.1"))
}

func TestWeightRoundRobinBalance(t *testing.T) {
	rb := &WeightRoundRobinBalance{}
	_ = rb.Add("127.0.0.1:2003", "4") //0
	_ = rb.Add("127.0.0.1:2004", "3") //1
	_ = rb.Add("127.0.0.1:2005", "2") //2

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
}
