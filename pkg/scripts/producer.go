package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/a1mart/kafkaesque/internal/generated/messaging" // Import your generated messaging service files
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run producer.go <topic>\n")
	}

	topic := os.Args[1]
	fmt.Printf("Running producer for topic: %s\n", topic)

	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Update if needed
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := messaging.NewMessagingServiceClient(conn)

	// Set up the input scanner to read messages
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to publish to the topic. Type 'exit' to stop.")

	// Loop to accept messages from the user
	for {
		fmt.Print("Message: ")
		scanner.Scan()
		messageText := scanner.Text()

		if messageText == "exit" {
			break
		}

		// Create message to publish
		message := &messaging.Message{
			Id:      fmt.Sprintf("%d", time.Now().UnixNano()),
			Type:    "test", // or any type you want to assign
			Payload: []byte(messageText),
		}
		publishReq := &messaging.PublishRequest{
			Topic:   topic,
			Message: message,
		}

		// Publish the message
		_, err := client.Publish(context.Background(), publishReq)
		if err != nil {
			log.Printf("Publish failed: %v", err)
		} else {
			fmt.Println("Message published!")
		}
	}
	fmt.Println("Producer stopped.")
}
