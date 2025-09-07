package mnemosyne

import (
	"sync"
	"time"
)

// MemTable structure
type MemTable struct {
	data map[string]interface{}
	mu   sync.RWMutex
	ttl  time.Duration
}

// NewMemTable initializes a MemTable with TTL
func NewMemTable(ttl time.Duration) *MemTable {
	memTable := &MemTable{
		data: make(map[string]interface{}),
		ttl:  ttl,
	}

	// Background cleanup routine
	go memTable.cleanupExpired()
	return memTable
}

// Put stores a message in the MemTable
func (mt *MemTable) Put(key string, value interface{}) {
	mt.mu.Lock()
	mt.data[key] = value
	mt.mu.Unlock()
}

// Get retrieves a message
func (mt *MemTable) Get(key string) (interface{}, bool) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	val, exists := mt.data[key]
	return val, exists
}

// Delete removes a message
func (mt *MemTable) Delete(key string) {
	mt.mu.Lock()
	delete(mt.data, key)
	mt.mu.Unlock()
}

// Cleanup expired messages
func (mt *MemTable) cleanupExpired() {
	for {
		time.Sleep(mt.ttl)
		mt.mu.Lock()
		for key := range mt.data {
			delete(mt.data, key)
		}
		mt.mu.Unlock()
	}
}
