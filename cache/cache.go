package cache

import (
	"log"
	"sync"
)

type Cache interface {
	Set(string,[]byte)error
	Get(string)([]byte,error)
	Del(string) error
	GetStat()  Stat
}


type Stat struct {
	Count int64
	KeySize int64
	ValueSize int64
}

func(s *Stat)add (k string, v []byte){
	s.Count ++
	s.KeySize +=int64(len(v))
	s.ValueSize += int64(len(v))
}

func(s *Stat)del( k string, v []byte){
	s.Count --
	s.KeySize-= int64(len(k))
	s.ValueSize -= int64(len(v))
}

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (cache *inMemoryCache) Set(k string,v []byte) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	tmp,exit := cache.c[k]
	if exit{
		cache.del(k,tmp)
	}

	cache.c[k] =v
	cache.add(k,v)
	return nil

}

func (cache *inMemoryCache) Get(k string) ([]byte, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	return cache.c[k],nil
}

func (cache *inMemoryCache) Del(k string) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	v, exit := cache.c[k]
	if exit{
		delete(cache.c,k)
		cache.del(k,v)
	}
	return nil
}

func (cache *inMemoryCache) GetStat() Stat {
	return cache.Stat
}


func New(typ string) Cache{
	var c Cache
	if typ =="inmemory"{
		c = newInmemoryCache()
	}
	if c == nil{
		panic("unknown cache type")
	}

	log.Println("SDSDSDSDSD")
	return c
}

func newInmemoryCache() *inMemoryCache{
	return &inMemoryCache{
		make(map[string][]byte),sync.RWMutex{},Stat{},
	}
}