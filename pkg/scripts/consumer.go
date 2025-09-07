// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/a1mart/kafkaesque/internal/generated/messaging" // Import your generated messaging service files
// 	"google.golang.org/grpc"
// )

// func main() {
// 	if len(os.Args) < 2 {
// 		log.Fatalf("Usage: go run consumer.go <topic>\n")
// 	}

// 	topic := os.Args[1]
// 	fmt.Printf("Running consumer for topic: %s\n", topic)

// 	// Connect to the gRPC server
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Update if needed
// 	if err != nil {
// 		log.Fatalf("failed to connect: %v", err)
// 	}
// 	defer conn.Close()

// 	client := messaging.NewMessagingServiceClient(conn)

// 	// Set up a signal channel to gracefully shut down
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// 	// Initialize consumer group ID and last processed message cursor (you can use dynamic values or more persistent mechanisms in real-world scenarios)
// 	consumerGroup := 1
// 	lastProcessedMessage := int64(0) // Start at 0 or adjust based on your implementation

// 	// Start consuming messages in a loop
// 	go func() {
// 		for {
// 			// You might want to adjust the batch size or other parameters dynamically here
// 			consumeReq := &messaging.ConsumeRequest{
// 				Topic:         topic,
// 				ConsumerGroup: int32(consumerGroup), // Group ID
// 				BatchSize:     1,                    // You can change this based on your needs
// 			}

// 			// Request messages
// 			resp, err := client.Consume(context.Background(), consumeReq)
// 			if err != nil {
// 				log.Printf("Consume failed: %v", err)
// 				time.Sleep(1 * time.Second) // Adjust delay based on backoff policy
// 				continue
// 			}

// 			// Process consumed messages
// 			if len(resp.GetMessages()) == 0 {
// 				log.Println("No new messages available. Retrying...")
// 			}

// 			for _, msg := range resp.GetMessages() {
// 				fmt.Printf("Consumed message: %s\n", msg.GetPayload())

// 				// Update lastProcessedMessage if necessary
// 				lastProcessedMessage++ // Assuming you can track the last processed message ID in the sequence
// 			}

// 			// Add a delay or backoff here if you're consuming at a high rate
// 			time.Sleep(500 * time.Millisecond) // Adjust based on how fast you want the consumer to poll
// 		}
// 	}()

// 	// Wait for a termination signal to shut down gracefully
// 	<-sigChan
// 	fmt.Println("Consumer stopped.")
// }

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/a1mart/kafkaesque/internal/generated/messaging" // Import your generated messaging service files
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run consumer.go <topic>\n")
	}

	topic := os.Args[1]
	fmt.Printf("Running consumer for topic: %s\n", topic)

	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Update if needed
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := messaging.NewMessagingServiceClient(conn)

	// Set up a signal channel to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages in a loop
	go func() {
		for {
			consumeReq := &messaging.ConsumeRequest{
				Topic:         topic,
				ConsumerGroup: 1, // Could make this dynamic if needed
				BatchSize:     2,
			}
			resp, err := client.Consume(context.Background(), consumeReq)
			if err != nil {
				log.Printf("Consume failed: %v", err)
			} else {
				for _, msg := range resp.GetMessages() {
					fmt.Printf("Consumed message: %s\n", msg.GetPayload())
				}
			}
		}
	}()

	// Wait for a termination signal to shut down gracefully
	<-sigChan
	fmt.Println("Consumer stopped.")
}
