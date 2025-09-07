package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/a1mart/kafkaesque/internal/generated/messaging"
	"google.golang.org/grpc"
)

const serverAddr = "localhost:50051"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run pkg/scripts/create_topic.go <topic_name> <strategy>")
		os.Exit(1)
	}

	topicName := os.Args[1]
	strategy := os.Args[2]

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := messaging.NewAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &messaging.CreateTopicRequest{
		Topic:    topicName,
		Strategy: strategy,
	}

	resp, err := client.CreateTopic(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	if resp.Success {
		fmt.Printf("Successfully created topic: %s with strategy: %s\n", topicName, strategy)
	} else {
		fmt.Printf("Failed to create topic: %s. Error: %s\n", topicName, resp.Error)
	}
}
