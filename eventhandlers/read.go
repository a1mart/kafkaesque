package eventhandlers

import (
	"bufio"
	"fmt"
	"os"
)

// ReadMessagesFromFile reads messages from the specified file
func ReadMessagesFromFile(filePath string) ([]string, error) {
	var messages []string

	// Open the file in read-only mode
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read each line from the file and store in messages slice
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		messages = append(messages, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from file: %v", err)
	}

	return messages, nil
}

// // Example usage inside the consumer loop
// func main() {
// 	// Set the file path for reading stored messages
// 	filePath := "messages.log"

// 	// Read messages from file
// 	messages, err := ReadMessagesFromFile(filePath)
// 	if err != nil {
// 		log.Fatalf("Failed to read messages from file: %v", err)
// 	}

// 	// Print the messages read from the file
// 	fmt.Println("Messages read from file:")
// 	for _, msg := range messages {
// 		fmt.Println(msg)
// 	}

// 	// Consumer logic would continue here, for example, by consuming new messages from the gRPC service
// }
