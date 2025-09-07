package server

import (
	"context"
	"log"

	"github.com/a1mart/kafkaesque/internal/generated/messaging"
)

// Admin Service: Create Topic
func (s *Server) CreateTopic(ctx context.Context, req *messaging.CreateTopicRequest) (*messaging.CreateTopicResponse, error) {
	topic := req.GetTopic()
	strategy := req.GetStrategy()

	if topic == "" || strategy == "" {
		return &messaging.CreateTopicResponse{Success: false, Error: "Invalid topic or strategy"}, nil
	}

	if _, exists := s.topics[topic]; exists {
		return &messaging.CreateTopicResponse{Success: false, Error: "Topic already exists"}, nil
	}

	s.topics[topic] = strategy
	log.Printf("Created topic: %s with strategy: %s", topic, strategy)
	return &messaging.CreateTopicResponse{Success: true}, nil
}

// Admin Service: List Topics
func (s *Server) ListTopics(ctx context.Context, req *messaging.ListTopicsRequest) (*messaging.ListTopicsResponse, error) {
	var topicList []*messaging.TopicInfo

	for topic, strategy := range s.topics {
		topicList = append(topicList, &messaging.TopicInfo{Topic: topic, Strategy: strategy})
	}

	return &messaging.ListTopicsResponse{Topics: topicList}, nil
}

func (s *Server) ListConsumers(ctx context.Context, req *messaging.ListConsumersRequest) (*messaging.ListConsumersResponse, error) {
	var consumers []string
	for c := range s.consumerMap {
		consumers = append(consumers, c)
	}
	return &messaging.ListConsumersResponse{ConsumerGroups: consumers}, nil
}
