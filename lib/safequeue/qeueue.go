package safequeue

import (
	"sync"
	"errors"
)

//协程安全.
type SafeQueue struct {
	q []string
	w sync.RWMutex
}

//入队时加写独占锁
func (u *SafeQueue) Push(val string) {
	u.w.Lock()
	defer u.w.Unlock()

	u.q = append(u.q, val)
}

//出队时加读共享锁
func (u *SafeQueue) Pop() (string, error) {
	u.w.Lock()
	defer u.w.Unlock()

    if len(u.q) > 0 {
        ret := u.q[0]
        u.q = u.q[1:]

        return ret, nil
    }

    return "", errors.New("queue empty")
}

//查看队列长度
func (u *SafeQueue) Len() int {
    return len(u.q)
}

//初始化一个map
func Create() SafeQueue {
	return SafeQueue{w: sync.RWMutex{}, q: []string{}}
}
