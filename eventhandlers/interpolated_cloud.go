package eventhandlers

import (
	"fmt"
	"sync"
)

/*
unique key, latest value
	history of values for a key
	transaction ~ atomic... i.e. add multiple KVP in an operation
	nesting as an epoch
*/

// ICloud is a struct that stores key-value pairs and their histories
type ICloud struct {
	data       map[string][]string   // Store history as slices of values
	mu         sync.RWMutex          // Mutex for thread-safe operations
	epochStack []map[string][]string // Stack to store nested epoch transactions (maps)
}

// NewICloud creates and returns a new ICloud instance
func NewICloud() *ICloud {
	return &ICloud{
		data:       make(map[string][]string),
		epochStack: []map[string][]string{}, // Initialize empty stack for epochs
	}
}

// AddOrUpdate adds a new value to the history of a key
func (cloud *ICloud) AddOrUpdate(key, value string) {
	cloud.mu.Lock()
	defer cloud.mu.Unlock()

	cloud.data[key] = append(cloud.data[key], value) // Append value to history
	fmt.Printf("Added/Updated key: %s with value: %s\n", key, value)
}

// GetLatest retrieves the latest value for a given key
func (cloud *ICloud) GetLatest(key string) (string, bool) {
	cloud.mu.RLock()
	defer cloud.mu.RUnlock()

	history, exists := cloud.data[key]
	if !exists || len(history) == 0 {
		return "", false
	}

	// Return the latest value (last element in the history)
	return history[len(history)-1], true
}

// GetHistory retrieves the entire history of values for a given key
func (cloud *ICloud) GetHistory(key string) ([]string, bool) {
	cloud.mu.RLock()
	defer cloud.mu.RUnlock()

	history, exists := cloud.data[key]
	return history, exists
}

// StartEpoch begins a new epoch (group of transactions)
func (cloud *ICloud) StartEpoch() {
	cloud.mu.Lock()
	defer cloud.mu.Unlock()

	// Push the current data to the epoch stack to save it
	epochCopy := make(map[string][]string)
	for k, v := range cloud.data {
		epochCopy[k] = append([]string(nil), v...) // Copy the current history
	}

	// Push the copy of the current data into the stack
	cloud.epochStack = append(cloud.epochStack, epochCopy)
	fmt.Println("Started new epoch")
}

// CommitEpoch commits the current epoch and applies the changes
func (cloud *ICloud) CommitEpoch() {
	cloud.mu.Lock()
	defer cloud.mu.Unlock()

	// Pop the top epoch from the stack (commit the changes)
	if len(cloud.epochStack) == 0 {
		fmt.Println("No epoch to commit")
		return
	}
	// For now, we just print instead of popping it immediately
	fmt.Println("Committed current epoch. Keeping epoch in stack.")
}

// RollbackEpoch rolls back the current epoch (undo changes made during the epoch)
func (cloud *ICloud) RollbackEpoch() {
	cloud.mu.Lock()
	defer cloud.mu.Unlock()

	// Roll back the changes to the last saved epoch
	if len(cloud.epochStack) == 0 {
		fmt.Println("No epoch to rollback")
		return
	}

	// Pop the top epoch from the stack
	epochCopy := cloud.epochStack[len(cloud.epochStack)-1]
	cloud.data = epochCopy
	cloud.epochStack = cloud.epochStack[:len(cloud.epochStack)-1]

	fmt.Println("Rolled back to previous epoch")
}

// AddOrUpdateBatch performs a batch update for multiple key-value pairs atomically
func (cloud *ICloud) AddOrUpdateBatch(updates map[string]string) {
	cloud.mu.Lock()
	defer cloud.mu.Unlock()

	// Perform the batch update atomically
	for key, value := range updates {
		cloud.data[key] = append(cloud.data[key], value)
		fmt.Printf("Added/Updated key: %s with value: %s in batch\n", key, value)
	}
}

func main() {
	// Creating an instance of ICloud
	cloud := NewICloud()

	// Adding/updating some keys
	cloud.AddOrUpdate("key1", "value1")
	cloud.AddOrUpdate("key2", "value2")
	cloud.AddOrUpdate("key3", "value3")
	cloud.AddOrUpdate("key3", "value4")
	cloud.AddOrUpdate("key1", "value5")

	// Start a new epoch
	cloud.StartEpoch()

	// Perform batch update inside the epoch
	cloud.AddOrUpdateBatch(map[string]string{
		"key4": "value6",
		"key5": "value7",
	})

	// Retrieve the latest value for a key
	if value, exists := cloud.GetLatest("key1"); exists {
		fmt.Println("Latest value for key1:", value)
	} else {
		fmt.Println("key1 does not exist")
	}

	// Commit the current epoch
	cloud.CommitEpoch()

	// Retrieve the entire history for a key
	if history, exists := cloud.GetHistory("key3"); exists {
		fmt.Println("History for key3:", history)
	} else {
		fmt.Println("key3 does not exist")
	}

	// Get all key-value pairs (just the latest values)
	allData := make(map[string]string)
	for key, history := range cloud.data {
		allData[key] = history[len(history)-1] // Store only the latest value in the map
	}
	fmt.Println("All key-value pairs:", allData)

	// Rollback epoch and print values
	cloud.RollbackEpoch()
	allDataAfterRollback := make(map[string]string)
	for key, history := range cloud.data {
		allDataAfterRollback[key] = history[len(history)-1]
	}
	fmt.Println("All key-value pairs after rollback:", allDataAfterRollback)
}
