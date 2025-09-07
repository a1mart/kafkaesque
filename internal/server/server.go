package server

import (
	"time"

	"github.com/a1mart/kafkaesque/internal/draupnir"
	"github.com/a1mart/kafkaesque/internal/generated/messaging"
	"github.com/a1mart/kafkaesque/internal/mnemosyne"
)

// gRPC Server
type Server struct {
	messaging.UnimplementedMessagingServiceServer
	messaging.UnimplementedAdminServiceServer // Implement AdminService
	rb                                        *draupnir.RingBuffer
	memTable                                  *mnemosyne.MemTable
	topics                                    map[string]string // Store topics and their strategies
	consumerMap                               map[string]bool   // Tracks registered consumers
}

func NewServer(size int, numConsumers int, ttl time.Duration) *Server {
	s := &Server{
		rb:          draupnir.NewRingBuffer(size, numConsumers),
		memTable:    mnemosyne.NewMemTable(ttl),
		topics:      make(map[string]string),
		consumerMap: make(map[string]bool),
	}
	return s
}
