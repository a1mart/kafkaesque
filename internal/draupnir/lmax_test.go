package draupnir_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/a1mart/kafkaesque/internal/draupnir"
	"github.com/a1mart/kafkaesque/internal/generated/messaging"
)

// TestNewRingBuffer ensures the buffer initializes correctly.
func TestNewRingBuffer(t *testing.T) {
	size := 5
	numConsumers := 2
	rb := draupnir.NewRingBuffer(size, numConsumers)

	if rb == nil {
		t.Fatal("RingBuffer should not be nil")
	}
	if len(rb.ReadCursors) != numConsumers {
		t.Errorf("Expected %d consumers, got %d", numConsumers, len(rb.ReadCursors))
	}
}

// TestPutAndGet ensures messages are written and read correctly.
func TestPutAndGet(t *testing.T) {
	rb := draupnir.NewRingBuffer(3, 1)

	msg1 := &messaging.Message{Id: "1"}
	msg2 := &messaging.Message{Id: "2"}

	rb.Put(msg1, msg2)

	msgs := rb.Get(2, 0)

	if len(msgs) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Id != "1" || msgs[1].Id != "2" {
		t.Error("Messages not retrieved in correct order")
	}
}

// TestRingBufferWrapAround ensures data wraps correctly when buffer is full.
func TestRingBufferWrapAround(t *testing.T) {
	rb := draupnir.NewRingBuffer(3, 1)

	msg1 := &messaging.Message{Id: "1"}
	msg2 := &messaging.Message{Id: "2"}
	msg3 := &messaging.Message{Id: "3"}
	msg4 := &messaging.Message{Id: "4"}

	rb.Put(msg1, msg2, msg3)
	rb.Get(2, 0)
	rb.Put(msg4)

	msgs := rb.Get(2, 0)
	if len(msgs) != 2 || msgs[0].Id != "3" || msgs[1].Id != "4" {
		t.Error("Ring buffer wrap-around failed")
	}
}

// TestConcurrentAccess ensures thread safety of Put and Get.
func TestConcurrentAccess(t *testing.T) {
	rb := draupnir.NewRingBuffer(10, 2)
	var wg sync.WaitGroup

	producer := func(id string) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			msg := &messaging.Message{Id: id}
			rb.Put(msg)
		}
	}

	consumer := func(group int) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			rb.Get(1, group)
		}
	}

	wg.Add(4)
	go producer("A")
	go producer("B")
	go consumer(0)
	go consumer(1)
	wg.Wait()
}

// TestAtomicCounters ensures atomic operations work as expected.
func TestAtomicCounters(t *testing.T) {
	_ = draupnir.NewRingBuffer(5, 1)
	var writeCursor int64
	atomic.StoreInt64(&writeCursor, 2)
	atomic.AddInt64(&writeCursor, 1)
	if atomic.LoadInt64(&writeCursor) != 3 {
		t.Error("Atomic counter did not increment correctly")
	}
}
