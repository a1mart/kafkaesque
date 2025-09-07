package eventhandlers

import (
	"fmt"
	"log"
	"os"
)

// WriteMessageToFile writes the message to a specified file
func WriteMessageToFile(message string, filePath string) error {
	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Write the message to the file
	_, err = file.WriteString(message + "\n")
	if err != nil {
		return fmt.Errorf("failed to write message to file: %v", err)
	}
	log.Printf("Message written to file: %s", message)
	return nil
}

// // Example usage inside the producer loop
// func main() {
// 	// Set the file path for storing messages
// 	filePath := "messages.log"

// 	// Producer logic
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("Enter messages to publish to the topic. Type 'exit' to stop.")

// 	for {
// 		fmt.Print("Message: ")
// 		scanner.Scan()
// 		messageText := scanner.Text()

// 		if messageText == "exit" {
// 			break
// 		}

// 		// Create message to publish
// 		message := &messaging.Message{
// 			Id:      fmt.Sprintf("%d", time.Now().UnixNano()),
// 			Type:    "test",
// 			Payload: []byte(messageText),
// 		}
// 		publishReq := &messaging.PublishRequest{
// 			Topic:   "test-topic",
// 			Message: message,
// 		}

// 		// Publish the message
// 		_, err := client.Publish(context.Background(), publishReq)
// 		if err != nil {
// 			log.Printf("Publish failed: %v", err)
// 		} else {
// 			fmt.Println("Message published!")

// 			// Write message to file after publishing
// 			err := WriteMessageToFile(messageText, filePath)
// 			if err != nil {
// 				log.Printf("Failed to write message to file: %v", err)
// 			}
// 		}
// 	}
// 	fmt.Println("Producer stopped.")
// }
