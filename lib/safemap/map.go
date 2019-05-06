package safemap

import (
	"sync"
)

//map加锁,协程安全.
type SafeMap struct {
	m map[string]bool
	w sync.RWMutex
}

//设置时加写独占锁
func (u *SafeMap) Set(key string, val bool) {
	u.w.Lock()
	defer u.w.Unlock()

	u.m[key] = val
}

//获取时加读共享锁
func (u *SafeMap) Get(key string) bool {
	u.w.RLock()
	defer u.w.RUnlock()

	ret, ok := u.m[key]

	if ok {
		return ret
	}

	return false
}

//初始化一个map
func Create() SafeMap {
	return SafeMap{w: sync.RWMutex{}, m: make(map[string]bool)}
}
