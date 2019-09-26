package main

import (
	"log"
	"sync"
	"time"
)

type cacheEntry struct {
	v  []byte
	ts int64
}

var m map[string]cacheEntry
var l sync.RWMutex

func init() {
	m = make(map[string]cacheEntry)
	go cleaner()
	log.Print("[I] Cache module loaded.")
}

func cleaner() {
	for range time.Tick(60 * time.Second) {
		l.Lock()
		for k, v := range m {
			if v.ts > 0 && v.ts < getNowMillisecond() {
				// Expired
				delete(m, k)
			}
		}
		l.Unlock()
	}
}

func cacheSet(k string, v []byte, ttlMs int) {
	var t int64
	if ttlMs > 0 {
		t = getNowMillisecond() + int64(ttlMs)
	} else {
		t = -1
	}
	l.Lock()
	m[k] = cacheEntry{v: v, ts: t}
	l.Unlock()
}

func cacheGet(k string) []byte {
	l.RLock()
	e, ok := m[k]
	l.RUnlock()
	if ok && e.ts > 0 && getNowMillisecond() <= e.ts {
		return e.v
	}
	return nil
}

func getNowMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
