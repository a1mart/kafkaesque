// connectors/sources/kafka/kafka_source.go
package kafka

import (
	"fmt"

	connectors "github.com/a1mart/kafkaesque/internal/indranet"
)

type KafkaSource struct {
	brokers []string
	topic   string
}

func (k *KafkaSource) Init(config map[string]string) error {
	if brokers, ok := config["brokers"]; ok {
		k.brokers = []string{brokers}
	} else {
		return fmt.Errorf("missing 'brokers' config")
	}

	if topic, ok := config["topic"]; ok {
		k.topic = topic
	} else {
		return fmt.Errorf("missing 'topic' config")
	}

	fmt.Println("Kafka source initialized")
	return nil
}

func (k *KafkaSource) Start() error {
	fmt.Printf("Starting Kafka source: brokers=%v, topic=%s\n", k.brokers, k.topic)
	return nil
}

func (k *KafkaSource) Stop() error {
	fmt.Println("Stopping Kafka source")
	return nil
}

// Factory function
func NewKafkaSource() connectors.Connector {
	return &KafkaSource{}
}
