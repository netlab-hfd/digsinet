package iface

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lachstec/digsinet-ng/config"
	"github.com/Lachstec/digsinet-ng/event"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api"

	"google.golang.org/protobuf/encoding/prototext"
)

// GNMIHandler handles gNMI subscriptions and publishes notifications to Kafka.
type GNMIHandler struct {
	*event.KafkaHandler
}

// NewGNMIHandler creates a new GNMIHandler.
func NewGNMIHandler() (*GNMIHandler, error) {
	kh, err := event.NewKafkaHandler()
	return &GNMIHandler{
		kh,
	}, err
}

// SubscribeAndPublish subscribes to gNMI paths and publishes notifications to Kafka.
func (h *GNMIHandler) SubscribeAndPublish(address string, paths []string, topic string) error {
	conf := config.GetConfig()

	// create a gNMI Target erstellen
	tg, err := api.NewTarget(
		api.Name("gnmi_node1"),       // TODO: make configurable, or remove
		api.Address(address+":6030"), // needed for arista_ceos TODO: make configurable, default to 6030
		api.Username(conf.GetString("gnmi.username")),
		api.Password(conf.GetString("gnmi.password")),
		api.SkipVerify(true), // TODO: make configurable, default to false
		api.Insecure(true),   // needed for arista_ceos TODO: make configurable, default to false
	)
	if err != nil {
		return fmt.Errorf("failed to create gNMI target: %v", err)
	}

	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create a gNMI client
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create gNMI client: %v", err)
	}
	defer tg.Close()

	// send a gNMI capabilities request to the created target
	capResp, err := tg.Capabilities(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prototext.Format(capResp))

	if conf.GetBool("gnmi.publish") {
		// publish dummy notification to Kafka
		err = h.KafkaHandler.PublishGNMINotificationToKafka(topic, &gnmi.Path{}, &gnmi.TypedValue{})
		if err != nil {
			return fmt.Errorf("failed to publish gNMI notification to Kafka: %v", err)
		}
	}

	return nil
}
