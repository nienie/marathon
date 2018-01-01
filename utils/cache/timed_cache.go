package cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	defaultCleanupInterval               = 100 * time.Millisecond // 100ms
	noExpireTime           time.Duration = -1
	noExpireTimeFlag       time.Duration = 0
)

var (
	errKeyNotExist      = fmt.Errorf("key does not exist")
	errNoOnLoadCallback = fmt.Errorf("no OnLoad Callback")
)

//Callback ...
type Callback interface {
	//OnLoad it will be called when get a not existed key.
	OnLoad(key interface{}) (interface{}, error)

	//OnRemove it will be called before a cache is deleted
	OnRemove(key interface{}, val interface{}) error
}

//Item ...
type Item struct {
	Object             interface{}
	Expiration         *time.Time
	ExpirationInterval time.Duration
}

//Expired ...
func (item *Item) Expired() bool {
	if item.Expiration == nil {
		return false
	}
	return item.Expiration.Before(time.Now())
}

//TimedCache ...
type TimedCache struct {
	sync.RWMutex
	items             map[interface{}]*Item
	janitor           *janitor
	callback          Callback
	defaultExpireTime time.Duration
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor)Run(c *TimedCache) {
	j.stop = make(chan bool)
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.DelExpiredKeys()
		case <-j.stop:
			c.Clear()
			return
		}
	}
}

func (j *janitor)Stop() {
	j.stop <- true
}

func stopJanitor(c *TimedCache) {
	c.janitor.Stop()
}

func runJanitor(c *TimedCache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
	}
	c.janitor = j
	go j.Run(c)
}

//NewTimedCache ...
func NewTimedCache(defaultExpireTime time.Duration, callback Callback) *TimedCache {
	if defaultExpireTime <= noExpireTimeFlag {
		defaultExpireTime = noExpireTime
	}
	C := &TimedCache{
		items:             make(map[interface{}]*Item),
		callback:          callback,
		defaultExpireTime: defaultExpireTime,
	}
	runJanitor(C, defaultCleanupInterval)
	runtime.SetFinalizer(C, stopJanitor)
	return C
}

//Set ...
func (c *TimedCache)Set(key interface{}, val interface{}, expireTime time.Duration) error {
	c.Lock()
	defer c.Unlock()
	c.del(key)
	return c.set(key, val, expireTime)
}

func (c *TimedCache)set(key interface{}, val interface{}, expireTime time.Duration) error {
	item := &Item{
		Object: val,
	}
	if expireTime == noExpireTimeFlag {
		expireTime = c.defaultExpireTime
	}
	if expireTime > noExpireTimeFlag {
		expiration := time.Now().Add(expireTime)
		item.Expiration = &expiration
		item.ExpirationInterval = expireTime
	}
	c.items[key] = item
	return nil
}

//Del ...
func (c *TimedCache)Del(key interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.del(key)
}

func (c *TimedCache)del(key interface{}) error {
	var err error
	item, ok := c.items[key]
	if !ok {
		return nil
	}
	if c.callback != nil {
		err = c.callback.OnRemove(key, item.Object)
	}
	delete(c.items, key)
	return err
}

//Get ...
func (c *TimedCache)Get(key interface{}) (interface{}, error) {
	c.RLock()
	val, err := c.get(key)
	c.RUnlock()
	return val, err
}

func (c *TimedCache)get(key interface{}) (interface{}, error) {
	item, found := c.items[key]
	if !found || item.Expired() {
		return nil, errKeyNotExist
	}
	return item.Object, nil
}

//GetAndSetWhenNotExisted ...
func (c *TimedCache)GetAndSetWhenNotExisted(key interface{}) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	val, err := c.getAndRefreshExpireTime(key)
	if err == nil {
		return val, nil
	}
	if c.callback == nil {
		return nil, errNoOnLoadCallback
	}
	newVal, err := c.callback.OnLoad(key)
	if err != nil {
		return nil, err
	}
	c.set(key, newVal, time.Duration(0))
	return newVal, nil
}

//GetAndRefreshExpireTime ...
func (c *TimedCache)GetAndRefreshExpireTime(key interface{}) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	return c.getAndRefreshExpireTime(key)
}

func (c *TimedCache)getAndRefreshExpireTime(key interface{}) (interface{}, error) {
	item, found := c.items[key]
	if !found || item.Expired() {
		return nil, errKeyNotExist
	}
	if item.Expiration != nil && item.ExpirationInterval > 0 {
		newExpiration := item.Expiration.Add(item.ExpirationInterval)
		item.Expiration = &newExpiration
	}
	return item.Object, nil
}

//DelExpiredKeys ...
func (c *TimedCache)DelExpiredKeys() {
	c.Lock()
	defer c.Unlock()
	for key, item := range c.items {
		if item.Expired() {
			if c.callback != nil {
				c.callback.OnRemove(key, item.Object)
			}
			delete(c.items, key)
		}
	}
}

//ToMap ...
func (c *TimedCache)ToMap() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	c.RLock()
	defer c.RUnlock()
	for key, item := range c.items {
		if !item.Expired() {
			m[key] = item.Object
		}
	}
	return m
}

//Clear ...
func (c *TimedCache)Clear() {
	c.Lock()
	defer c.Unlock()
	for key, item := range c.items {
		if c.callback != nil {
			c.callback.OnRemove(key, item.Object)
		}
	}
}