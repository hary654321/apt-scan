package pool

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestPool_NewTick(t *testing.T) {
	startTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	fmt.Println(strconv.FormatInt(rand.Int63(), 10))
	elapsed := time.Since(startTime)
	fmt.Println(elapsed)
}

func TestParse(t *testing.T) {
	println(111)
	fmt.Println(url.Parse("http://example.com/x/y%2Fz"))
}

// var res map[int]string
func TestPool(t *testing.T) {
	spyPool := New(9)

	res := make(map[int]string)

	spyPool.Function = func(i interface{}) {
		println("干活了", i.(int))

		res[i.(int)] = "aa"
		time.Sleep(1 * time.Second)
	}
	i := 1
	go func() {
		for i < 10 {
			i++
			spyPool.Push(i)
		}
		spyPool.Stop()
	}()

	spyPool.Run()

	fmt.Printf("%v", res)

}
