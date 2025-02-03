package event

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Lachstec/digsinet-ng/config"
	"github.com/openconfig/gnmi/proto/gnmi"
)

type KafkaHandler struct {
	producer sarama.SyncProducer
}

func NewKafkaHandler() (*KafkaHandler, error) {
	configKafka := sarama.NewConfig()
	configKafka.Producer.Return.Successes = true

	conf := config.GetConfig()
	kafkaBrokers := conf.GetStringSlice("kafka.brokers")

	producer, err := sarama.NewSyncProducer(kafkaBrokers, configKafka)
	return &KafkaHandler{
		producer,
	}, err
}

// PublishToKafka publishes the gNMI update to Kafka.
func (k KafkaHandler) PublishGNMINotificationToKafka(topic string, path *gnmi.Path, val *gnmi.TypedValue) error {
	message := fmt.Sprintf("gNMI Notification Path: %s, Value: %s", path, val)
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})
	return err
}
