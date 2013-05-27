package registry

import "sync"

// Registry holds references to last position for given file
// and it needs to be updated whenever a file is read...
// it needs locks and such because it will be accessed by couple of routines
type Registry struct {
	lock  sync.RWMutex
	store map[string]int64
}

func (r *Registry) Get(key string) int64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	val := r.store[key]
	return val
}

func (r *Registry) Set(key string, val int64) int64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.store[key] = val
	return val

}

func New() *Registry {
	return &Registry{store: make(map[string]int64)}
}
