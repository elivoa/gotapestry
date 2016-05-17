package service

import (
	"fmt"
	"strconv"
	"syd/dal/constdao"
	"sync"
	"time"
)

// support models
type ConstCache struct {
	// name -> key -> ConstItem
	data map[string]map[string]*ConstCacheItem
	l    sync.RWMutex
}

type ConstCacheItem struct {
	Value      string
	FloatValue float64
	timeout    int64 // time
}

// services

func NewCosntService() *ConstService {
	return &ConstService{
		data: &ConstCache{
			data: map[string]map[string]*ConstCacheItem{},
		},
	}
}

type ConstService struct {
	data *ConstCache
}

// Cache系统，如果没有值，则每次都要读取数据库。
// 如果有值，则只读取一次，TODO 加上一个超时时间。重新读取数据。

func (s *ConstService) GetIntValue(name, key string) (int64, error) {
	ccitem := s.get(name, key)
	return strconv.ParseInt(ccitem.Value, 10, 64)
}

func (s *ConstService) GetStringValue(name, key string) (string, error) {
	ccitem := s.get(name, key)
	return ccitem.Value, nil
}

func (s *ConstService) Get2ndIntValue(name, key string) (int64, error) {
	ccitem := s.get(name, key)
	return int64(ccitem.FloatValue), nil
}

func (s *ConstService) Get2ndStringValue(name, key string) (string, error) {
	ccitem := s.get(name, key)
	return fmt.Sprint(ccitem.FloatValue), nil
}

// basic get system.

// get Item from cache first, if not found get from db and update cache.
func (s *ConstService) get(name, key string) *ConstCacheItem {
	// nil check
	if nil != s.data {
		s.data = &ConstCache{
			data: map[string]map[string]*ConstCacheItem{},
		}
	}

	// TODO Timeout check.

	s.data.l.RLock()
	defer s.data.l.RUnlock()

	// level one
	keymap, ok := s.data.data[name]
	if !ok || keymap == nil {
		keymap = map[string]*ConstCacheItem{}
	}

	// level two
	ccitem, ok := keymap[key]
	if !ok || ccitem == nil {
		// get value from database;
		if constmodel, err := constdao.GetOne(name, key); err != nil {
			panic(err)
		} else {
			ccitem = &ConstCacheItem{
				Value:      constmodel.Value,
				FloatValue: constmodel.FloatValue,
				timeout:    time.Now().UnixNano(),
			}
			keymap[key] = ccitem
		}
	}
	return ccitem
}

func (s *ConstService) Set(name string, key string, value interface{}, floatValue float64) error {
	if err := constdao.Set(name, key, value, floatValue); err != nil {
		return err
	}
	return s.updatecache(name, key, value, floatValue)
}

func (s *ConstService) Update(name string, key string, value interface{}, floatValue float64, id int64) error {
	if err := constdao.Update(name, key, value, floatValue, id); err != nil {
		return err
	}
	return s.updatecache(name, key, value, floatValue)
}

// update cache
func (s *ConstService) updatecache(name string, key string, value interface{}, floatValue float64) error {
	s.data.l.Lock()
	defer s.data.l.Unlock()

	// level one
	keymap, ok := s.data.data[name]
	if !ok || keymap == nil {
		keymap = map[string]*ConstCacheItem{}
	}
	keymap[key] = &ConstCacheItem{
		Value:      fmt.Sprint(value),
		FloatValue: floatValue,
		timeout:    time.Now().UnixNano(),
	}
	return nil
}

func (s *ConstService) deletecache(name string, key string) error {
	s.data.l.Lock()
	defer s.data.l.Unlock()

	// level one
	keymap, ok := s.data.data[name]
	if !ok || keymap == nil {
		return nil // 如果没发现namemap就不用删除了。
	}
	// 删除
	delete(keymap, key)
	return nil
}

func (s *ConstService) DeleteById(id int64) (int64, error) {
	var err error
	c, err := constdao.GetById(id)
	if err != nil {
		panic(err)
	}
	if aff, err := constdao.DeleteById(id); err != nil {
		return aff, err
	}
	// udpate cache
	err = s.deletecache(c.Name, c.Key)
	return 0, err
}

// TODO get list
