package main

import (
	"context"
	"fmt"
	"log"

	"github.com/a1mart/kafkaesque/internal/generated/messaging" // Import your generated messaging service files

	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := messaging.NewMessagingServiceClient(conn)

	// Publish a message
	message := &messaging.Message{
		Id:      "1",
		Type:    "test",
		Payload: []byte("Hello, World!"),
	}
	publishReq := &messaging.PublishRequest{
		Topic:   "example-topic",
		Message: message,
	}

	_, err = client.Publish(context.Background(), publishReq)
	if err != nil {
		log.Fatalf("Publish failed: %v", err)
	}
	fmt.Println("Message published!")

	// Consume messages
	consumeReq := &messaging.ConsumeRequest{
		Topic:         "example-topic",
		ConsumerGroup: 1,
		BatchSize:     2,
	}
	resp, err := client.Consume(context.Background(), consumeReq)
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	for _, msg := range resp.GetMessages() {
		fmt.Printf("Consumed message: %s\n", msg.GetPayload())
	}
}
