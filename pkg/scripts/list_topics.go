package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/a1mart/kafkaesque/internal/generated/messaging"
	"google.golang.org/grpc"
)

const serverAddr = "localhost:50051"

func main() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := messaging.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &messaging.ListTopicsRequest{}
	resp, err := client.ListTopics(ctx, req)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}

	fmt.Println("Existing Topics:")
	for _, topic := range resp.Topics {
		fmt.Printf("- %s (Strategy: %s)\n", topic.Topic, topic.Strategy)
	}
}
