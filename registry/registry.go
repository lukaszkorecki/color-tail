package registry

import "sync"

type AnyValue interface {}

// Registry is a thread-safe map with strings as keys and any value
// locking is required because underlying map will be read and written to
// many go routines
type Registry struct {
	lock  sync.RWMutex
	store map[string]AnyValue
}

func (r *Registry) Get(key string) (AnyValue, bool){
	r.lock.RLock()
	defer r.lock.RUnlock()
	val, status := r.store[key]
	return val, status
}

func (r *Registry) Set(key string, val AnyValue) AnyValue {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.store[key] = val
	return val

}

func New() *Registry {
	return &Registry{store: make(map[string]AnyValue)}
}
