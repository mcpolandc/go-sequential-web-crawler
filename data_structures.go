package main

import "sync"

// IThreadSafeMap - interface defining functions to be implemented
type IThreadSafeMap interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

// ThreadSafeMap - locking map struct
type ThreadSafeMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// Set - locking function to set item in map
func (tsm *ThreadSafeMap) Set(key string, value interface{}) {
	tsm.Lock()
	defer tsm.Unlock()

	tsm.items[key] = value
}

// Get - locking function to get value at given key in map
func (tsm *ThreadSafeMap) Get(key string) interface{} {
	tsm.Lock()
	defer tsm.Unlock()

	value, _ := tsm.items[key]

	return value
}
