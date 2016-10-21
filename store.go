package captcha

import (
	"github.com/qAison/captcha/cache"
)

//-----------------------------------------------------------------------------------------------------------//

type Store interface {
	Set(key, value string)
	Get(key string) (value string)
	GetDelete(key string) (value string)
}

//-----------------------------------------------------------------------------------------------------------//

type memoryStore struct {
	cache cache.Cache
}

func NewMemoryStore(ageTime int64) Store {
	return &memoryStore{
		cache: cache.New(ageTime),
	}
}

func (self *memoryStore) Set(key, value string) {
	self.cache.Set(key, value)
}

func (self *memoryStore) Get(key string) (value string) {
	if v := self.cache.Get(key); v != nil {
		value = v.(string)
	}
	return
}

func (self *memoryStore) GetDelete(key string) (value string) {
	if v := self.cache.Get(key); v != nil {
		value = v.(string)
		self.cache.Delete(key)
	}
	return
}

//-----------------------------------------------------------------------------------------------------------//
