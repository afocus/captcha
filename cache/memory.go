package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Delete(key interface{})
	UpdateTime(key interface{})
	Len() int
	Clear()
}

//-----------------------------------------------------------------------------------------------------------//

type memStore struct {
	key        interface{}
	value      interface{}
	createTime time.Time
}

type memCache struct {
	lock    sync.RWMutex
	values  map[interface{}]*list.Element
	list    *list.List
	ageTime int64
}

func New(ageTime int64) Cache {
	mc := &memCache{
		values: make(map[interface{}]*list.Element),
		list:   list.New(),
	}
	if ageTime <= 0 {
		ageTime = 60
	}
	mc.ageTime = ageTime
	go mc.gc()
	return mc
}

func (self *memCache) Len() int {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return len(self.values)
}

func (self *memCache) Set(key, value interface{}) {
	self.Delete(key)

	self.lock.Lock()
	defer self.lock.Unlock()

	val := memStore{
		key:        key,
		value:      value,
		createTime: time.Now(),
	}
	self.values[key] = self.list.PushFront(&val)
}

func (self *memCache) Get(key interface{}) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if el, ok := self.values[key]; ok {
		return el.Value.(*memStore).value
	} else {
		return nil
	}
}

func (self *memCache) UpdateTime(key interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if el, ok := self.values[key]; ok {
		el.Value.(*memStore).createTime = time.Now()
		self.list.MoveToFront(el)
	}
}

func (self *memCache) Delete(key interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if el, ok := self.values[key]; ok {
		delete(self.values, key)
		self.list.Remove(el)
	}
}

func (self *memCache) Clear() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.values = make(map[interface{}]*list.Element)
	self.list = list.New()
}

func (self *memCache) gc() {
	var interval int64
	for {
		self.lock.RLock()

		if el := self.list.Back(); el != nil {
			interval = el.Value.(*memStore).createTime.Unix() + self.ageTime - time.Now().Unix()
			if interval <= 0 {
				self.lock.RUnlock()
				self.lock.Lock()
				delete(self.values, el.Value.(*memStore).key)
				self.list.Remove(el)
				self.lock.Unlock()
				continue
			}
		} else {
			interval = 1
		}
		self.lock.RUnlock()
		time.Sleep(time.Second * time.Duration(interval))
	}
}
