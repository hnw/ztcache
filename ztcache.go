// Package ztcache は zero time cache pattern の実装です。
// 同じ結果を返す処理が並行に呼び出された場合に、1回の処理にまとめて相乗りさせるような仕組みを提供します。
package ztcache

import (
	//"log"
	"sync"
	"time"
)

// ZTCache は動作中の処理を管理する構造体です。
type ZTCache struct {
	m     sync.Mutex
	start map[string]time.Time
	chans map[string][]chan string
}

// New は ZTCache構造体のコンストラクタです。
func New() *ZTCache {
	return &ZTCache{
		start: make(map[string]time.Time),
		chans: make(map[string][]chan string),
	}
}

// Get は処理結果を取得するためのラッパー関数です。
// 同じkeyで同時に呼び出された処理がある場合に並行に走らせず、前の処理の終了を待った上で必要なら相乗りして実行結果を複数呼び出し間で共有します。
func (c *ZTCache) Get(key string, f func() string) (string, error) {
	c.m.Lock()
	if !c.start[key].IsZero() {
		// いま動作中のものの終了を待つ
		q := make(chan string)
		c.chans[key] = append(c.chans[key], q)
		c.m.Unlock()
		//log.Printf("Get(%s,f): waiting previous process...\n", key)
		<-q
	} else {
		c.m.Unlock()
	}

	c.m.Lock()
	if !c.start[key].IsZero() {
		//log.Printf("Get(%s,f): waiting current process...\n", key)
		q := make(chan string)
		c.chans[key] = append(c.chans[key], q)
		c.m.Unlock()
		return <-q, nil
	}
	c.start[key] = time.Now()
	c.m.Unlock()

	ret := f()

	c.m.Lock()
	chans := c.chans[key]
	delete(c.start, key)
	delete(c.chans, key)
	c.m.Unlock()

	for _, q := range chans {
		q <- ret
	}
	return ret, nil
}
