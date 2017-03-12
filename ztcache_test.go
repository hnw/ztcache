package ztcache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	m   sync.Mutex
	cnt int
)

func TestGet(t *testing.T) {
	var ret string
	cache := New()
	ret, _ = cache.Get("foo", func() string { return "foo" })
	if ret != "foo" {
		t.Fatal("doesn't match")
	}
}

func ExampleZtcache_Get() {
	m.Lock()
	cnt = 0
	m.Unlock()

	cache := New()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		ret, _ := cache.Get("foo", heavyFunc)
		fmt.Println(ret)
		wg.Done()
	}()

	time.Sleep(50 * time.Millisecond)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			ret, _ := cache.Get("foo", heavyFunc)
			fmt.Println(ret)
			wg.Done()
		}()
	}

	wg.Wait()

	// Output:
	// foo1
	// foo2
	// foo2
	// foo2
}

func heavyFunc() string {
	m.Lock()
	cnt++
	curCnt := cnt
	m.Unlock()
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("foo%d", curCnt)
}
