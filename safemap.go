package commons

import (
	"sync"
)

type safeMap struct {
	lock    *sync.RWMutex
	safeMap map[interface{}]interface{}
}

func NewSafeMap() *safeMap {
	return &safeMap{
		lock:    new(sync.RWMutex),
		safeMap: make(map[interface{}]interface{}),
	}
}

func (s *safeMap) Get(k interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if val, ok := s.safeMap[k]; ok {
		return val
	}
	return nil
}

func (s *safeMap) Set(k interface{}, v interface{}) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if val, ok := s.safeMap[k]; !ok {
		s.safeMap[k] = v
	} else if val != v {
		s.safeMap[k] = v
	} else {
		return false
	}
	return true
}

func (s *safeMap) Check(k interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.safeMap[k]
	return ok
}

func (s *safeMap) Delete(k interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.safeMap, k)
}

// Items returns all items in safemap.
func (s *safeMap) Items() map[interface{}]interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range s.safeMap {
		r[k] = v
	}
	return r
}

// Count returns the number of items within the map.
func (s *safeMap) Count() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.safeMap)
}
