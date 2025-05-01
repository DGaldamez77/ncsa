package main

import (
	"fmt"
	"sort"
	"time"

	cache "github.com/dgaldamez77/slidingcache"
)

func main() {
	testAddAndGet()

	testTreadSafety()

}

func testAddAndGet() {
	fmt.Println("starting at ", time.Now())
	c1 := cache.NewCache(nil)
	defer c1.Close()

	c1.Add("c1 item1", "value1", nil)
	c1.Add("c1 item2", "value2", duration(10))

	c2 := cache.NewCache(duration(15))
	defer c2.Close()

	c2.Add("c2 item1", "value1", duration(5))
	c2.Add("c2 item2", "value2", duration(13))

	time.Sleep(5 * time.Second)

	key := "c1 item2"
	v, err := c1.Get(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("value for key", key, "is", v)

	time.Sleep(70 * time.Second)

	key = "c1 item1"
	v, err = c1.Get(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("value for key", key, "is", v)
}
func testTreadSafety() {
	ch := make(chan int)
	done := make(chan struct{})
	c3 := cache.NewCache(nil)
	defer c3.Close()

	go addItems(ch, c3, 100)
	go removeItems(ch, done, c3, []int{4, 12, 26, 45, 55, 67, 98})

	time.Sleep(5 * time.Second)
	done <- struct{}{}
}

func duration(seconds int) *time.Duration {
	d := time.Duration(seconds * int(time.Second))
	return &d
}

func addItems(c chan<- int, cache *cache.Cache, n int) {
	for i := 0; i < n; i++ {
		cache.Add(fmt.Sprintf("item %d", i), fmt.Sprintf("value %d", i), nil)
		c <- i
	}
}

func removeItems(ch <-chan int, done chan struct{}, cache *cache.Cache, remove []int) {
	for {
		select {
		case i := <-ch:
			idx := sort.SearchInts(remove, i)
			if idx < len(remove) && remove[idx] == i {
				time.Sleep(100 * time.Millisecond)
				cache.Delete(fmt.Sprintf("item %d", i))
			}
		case <-done:
			return
		default:

		}
	}
}
