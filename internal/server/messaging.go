package server

import (
	"context"
	"fmt"
	"log"

	"github.com/a1mart/kafkaesque/internal/generated/messaging"
)

// Publish a message
func (s *Server) Publish(ctx context.Context, req *messaging.PublishRequest) (*messaging.PublishResponse, error) {
	message := req.GetMessage()
	if message == nil {
		log.Printf("Publish failed: message is nil")
		return nil, fmt.Errorf("invalid message")
	}

	s.rb.Put(message)
	log.Printf("Published message with ID: %s", message.GetId())
	return &messaging.PublishResponse{Success: true}, nil
}

// Consume messages
func (s *Server) Consume(ctx context.Context, req *messaging.ConsumeRequest) (*messaging.ConsumeResponse, error) {
	messages := s.rb.Get(int(req.GetBatchSize()), int(req.GetConsumerGroup()))
	if len(messages) == 0 {
		log.Printf("No messages to consume for consumer group %d", req.GetConsumerGroup())
	}

	for _, msg := range messages {
		if msg != nil {
			log.Printf("Consumed message with ID: %s", msg.GetId())
		}
	}

	return &messaging.ConsumeResponse{Messages: messages, Success: true}, nil
}

func (s *Server) Acknowledge(ctx context.Context, req *messaging.AckRequest) (*messaging.AckResponse, error) {
	log.Printf("Acknowledging messages: %v for consumer group %s", req.GetMessageIds(), req.GetConsumerGroup())
	return &messaging.AckResponse{Success: true}, nil
}

func (s *Server) GetDeadLetters(ctx context.Context, req *messaging.DeadLetterRequest) (*messaging.DeadLetterResponse, error) {
	log.Println("Fetching dead letter messages...")
	return &messaging.DeadLetterResponse{Messages: nil, Success: true}, nil
}

func (s *Server) RegisterConsumerGroup(ctx context.Context, req *messaging.RegisterConsumerRequest) (*messaging.RegisterConsumerResponse, error) {
	s.consumerMap[req.GetConsumerGroup()] = true
	log.Printf("Registered consumer group: %s", req.GetConsumerGroup())
	return &messaging.RegisterConsumerResponse{Success: true}, nil
}
