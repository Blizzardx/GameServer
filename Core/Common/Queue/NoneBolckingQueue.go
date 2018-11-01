package Queue

import (
	"sync"
)

// 不限制大小，添加不发生阻塞，接收阻塞等待
type NoneBlockingQueue struct {
	list      []interface{}
	listGuard sync.Mutex
	listCond  *sync.Cond
}

// 添加时不会发送阻塞
func (self *NoneBlockingQueue) Add(msg interface{}) {
	self.listGuard.Lock()
	self.list = append(self.list, msg)
	self.listGuard.Unlock()

	self.listCond.Signal()
}

func (self *NoneBlockingQueue) Reset() {
	self.list = self.list[0:0]
}

// 如果没有数据，发生阻塞
func (self *NoneBlockingQueue) Pick(retList *[]interface{}) (exit bool) {

	self.listGuard.Lock()

	for len(self.list) == 0 {
		self.listCond.Wait()
	}

	self.listGuard.Unlock()

	self.listGuard.Lock()

	// 复制出队列

	for _, data := range self.list {

		if data == nil {
			exit = true
			break
		} else {
			*retList = append(*retList, data)
		}
	}

	self.Reset()
	self.listGuard.Unlock()

	return
}

func NewNoneBlockingQueue() *NoneBlockingQueue {
	self := &NoneBlockingQueue{}
	self.listCond = sync.NewCond(&self.listGuard)

	return self
}
