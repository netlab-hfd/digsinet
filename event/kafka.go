package event

import (
	"github.com/IBM/sarama"
	"github.com/Lachstec/digsinet-ng/config"
)

type KafkaHandler struct {
	producer sarama.SyncProducer
}

func NewKafkaHandler(cfg config.Configuration) (*KafkaHandler, error) {
	configKafka := sarama.NewConfig()
	configKafka.Producer.Return.Successes = true

	var kafkaBrokers []string

	for _, broker := range cfg.Kafka.Brokers {
		kafkaBrokers = append(kafkaBrokers, broker.ConnectionString())
	}

	producer, err := sarama.NewSyncProducer(kafkaBrokers, configKafka)
	return &KafkaHandler{
		producer,
	}, err
}

// PublishToKafka publishes the gNMI update to Kafka.
func (k KafkaHandler) PublishGNMINotificationToKafka(topic string, msg string) error {
	//message := fmt.Sprintf("gNMI Notification Path: %s, Value: %s", msg)
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})
	return err
}
