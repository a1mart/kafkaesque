package draupnir

import (
	"log"
	"sync/atomic"

	"github.com/a1mart/kafkaesque/internal/generated/messaging"
)

// RingBuffer struct with advanced features
type RingBuffer struct {
	size         int
	buffer       []interface{}
	writeCursor  int64
	ReadCursors  []int64
	available    []int32 // Use int32 for availability tracking for atomic operations
	writeBarrier int64
	readBarrier  int64
}

// NewRingBuffer creates a new RingBuffer
func NewRingBuffer(size int, numConsumers int) *RingBuffer {
	rb := &RingBuffer{
		size:        size,
		buffer:      make([]interface{}, size),
		ReadCursors: make([]int64, numConsumers),
		available:   make([]int32, size), // Use int32 for availability tracking
	}
	for i := range rb.ReadCursors {
		rb.ReadCursors[i] = -1
	}
	return rb
}

// Get reads data from the ring buffer in a round-robin manner
func (rb *RingBuffer) Get(batchSize int, consumerGroup int) []*messaging.Message {
	results := make([]*messaging.Message, 0, batchSize)

	// Get the consumer's read cursor
	readCursor := &rb.ReadCursors[consumerGroup] // Track per-group cursor
	nextRead := atomic.LoadInt64(readCursor) + 1

	for i := 0; i < batchSize; {
		slot := nextRead % int64(rb.size) // Compute the ring buffer index

		if atomic.LoadInt32(&rb.available[slot]) == 1 {
			if msg, ok := rb.buffer[slot].(*messaging.Message); ok {
				results = append(results, msg)

				// Mark slot as consumed
				atomic.StoreInt32(&rb.available[slot], 0)
				atomic.StoreInt64(readCursor, nextRead) // Update cursor per group
				nextRead++
				i++ // Move to next batch item
			}
		} else {
			// If no available messages, break early to avoid unnecessary looping
			break
		}
	}

	return results
}

// // Get reads data from the ring buffer in batches (lock-free)
// func (rb *RingBuffer) Get(batchSize int, consumerId int) []*messaging.Message {
// 	readCursor := (*int64)(unsafe.Pointer(&rb.ReadCursors[consumerId])) // Unsafe cast for direct memory access
// 	nextRead := atomic.LoadInt64(readCursor) + 1
// 	results := make([]*messaging.Message, batchSize)

// 	index := 0 // Index for results slice
// 	for i := 0; i < batchSize; {
// 		if nextRead <= atomic.LoadInt64(&rb.writeCursor) {
// 			// Log the availability before consuming the message
// 			slot := nextRead % int64(rb.size)
// 			log.Printf("Consumer %d attempting to consume from Slot %d, Available: %d", consumerId, slot, atomic.LoadInt32(&rb.available[slot]))

// 			if atomic.LoadInt32(&rb.available[slot]) == 1 {
// 				if msg, ok := rb.buffer[slot].(*messaging.Message); ok {
// 					results[index] = msg
// 					log.Printf("Consumer %d consumed message ID: %s from Slot %d", consumerId, msg.GetId(), slot)

// 					// Mark the slot as unavailable after consumption
// 					atomic.StoreInt32(&rb.available[slot], 0)
// 					atomic.StoreInt64(readCursor, nextRead) // Update read cursor immediately after consumption
// 					nextRead++
// 					i++
// 					index++
// 				}
// 			} else {
// 				log.Printf("Consumer %d found Slot %d unavailable, retrying...", consumerId, slot)
// 				// Delay if the slot is unavailable, adjust as needed
// 				time.Sleep(1 * time.Microsecond)
// 			}

// 		} else {
// 			log.Printf("Consumer %d: No more messages to consume", consumerId)
// 			break
// 		}
// 	}

// 	// Log the final state of the read cursor
// 	log.Printf("Consumer %d finished consuming, readCursor: %d", consumerId, atomic.LoadInt64(readCursor))

// 	// Update the read barrier
// 	minReadCursor := atomic.LoadInt64(readCursor)
// 	for _, cursor := range rb.ReadCursors {
// 		if cursor < minReadCursor {
// 			minReadCursor = cursor
// 		}
// 	}
// 	atomic.StoreInt64(&rb.readBarrier, minReadCursor)

// 	return results[:index]
// }

// Put writes data to the ring buffer in batches
func (rb *RingBuffer) Put(data ...*messaging.Message) {
	nextWrite := atomic.AddInt64(&rb.writeCursor, int64(len(data))) - int64(len(data))
	for _, d := range data {
		slot := nextWrite % int64(rb.size)
		rb.buffer[slot] = d
		atomic.StoreInt32(&rb.available[slot], 1) // Mark as available
		log.Printf("Put Message ID: %s at Slot: %d, WriteCursor: %d, Available: %d", d.GetId(), slot, nextWrite, atomic.LoadInt32(&rb.available[slot]))
		nextWrite++
	}

	// Update the write cursor once all data is written
	atomic.StoreInt64(&rb.writeCursor, nextWrite)
	log.Printf("WriteCursor advanced to: %d", nextWrite)
}
