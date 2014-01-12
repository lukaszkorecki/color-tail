package registry

import "sync"

type AnyValue interface{}

// Registry is a thread-safe map with strings as keys and anything can be
// stored as a  value.
// Locking is required because underlying map will be read and written to
// by many go routines.
type Registry struct {
	lock  sync.RWMutex
	store map[string]AnyValue
}

// Returns value read from given key and true or false if the key exists
func (r *Registry) Get(key string) (AnyValue, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	val, status := r.store[key]
	return val, status
}

// Set (Safely) value under given key
func (r *Registry) Set(key string, val AnyValue) AnyValue {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.store[key] = val
	return val

}

func NewRegistry() *Registry {
	return &Registry{store: make(map[string]AnyValue)}
}
