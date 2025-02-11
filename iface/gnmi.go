package iface

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Lachstec/digsinet-ng/config"
	"github.com/Lachstec/digsinet-ng/event"

	"github.com/openconfig/gnmic/pkg/api"

	"google.golang.org/protobuf/encoding/prototext"
)

var nodeSubCtxCancelFuncs = make(map[string]map[string]context.CancelFunc)

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
func (gh *GNMIHandler) SubscribeAndPublish(address string, paths []string, target string) (string, error) {
	conf := config.GetConfig()

	// create a gNMI Target erstellen
	tg, err := api.NewTarget(
		api.Name(target),             // TODO: make configurable, or remove
		api.Address(address+":6030"), // needed for arista_ceos TODO: make configurable, default to 6030
		api.Username(conf.GetString("gnmi.username")),
		api.Password(conf.GetString("gnmi.password")),
		api.SkipVerify(true), // TODO: make configurable, default to false
		api.Insecure(true),   // needed for arista_ceos TODO: make configurable, default to false
	)
	if err != nil {
		return "", fmt.Errorf("failed to create gNMI target: %v", err)
	}

	// context with timeout
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	ctx := context.Background()

	// create a gNMI client
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create gNMI client: %v", err)
	}
	//defer tg.Close()

	// // send a gNMI capabilities request to the created target
	// capResp, err := tg.Capabilities(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(prototext.Format(capResp))

	// create a gNMI subscribeRequest
	var subscriptions []api.GNMIOption
	subscriptions = append(subscriptions,
		api.Encoding("json_ietf"),
		api.SubscriptionListMode("stream"))

	for _, path := range paths {
		subscriptions = append(subscriptions,
			api.Subscription(
				api.Path(path),
				api.SubscriptionMode("sample"),
				api.SampleInterval(10*time.Second),
			))
	}

	subReq, err := api.NewSubscribeRequest(subscriptions...)
	if err != nil {
		log.Fatal().AnErr("iface", err)
	}
	fmt.Println(prototext.Format(subReq))
	// Create a unique subscription ID
	subscriptionID := fmt.Sprintf("sub_%s_%d", address, time.Now().UnixNano())

	// Start the subscription in a goroutine
	subCtx, subCtxCancel := context.WithCancel(context.Background())
	// unnecessary?
	if len(nodeSubCtxCancelFuncs[target]) == 0 {
		nodeSubCtxCancelFuncs[target] = make(map[string]context.CancelFunc)
	}
	nodeSubCtxCancelFuncs[target][subscriptionID] = subCtxCancel
	go tg.Subscribe(subCtx, subReq, subscriptionID)

	// Start a goroutine to handle subscription cleanup
	go func() {
		<-subCtx.Done()
		log.Info().
			Str("Iface", "gNMI"). // should be Id() as for builder
			Msg("subscription cleanup done: " + subscriptionID)

		tg.StopSubscription(subscriptionID)
		tg.Close()
	}()

	// Create error channel for this subscription
	errChan := make(chan error, 1)

	// Handle subscription responses in a separate goroutine
	go func() {
		subRspChan, subErrChan := tg.ReadSubscriptions()
		for {
			select {
			case rsp := <-subRspChan:
				log.Debug().
					Str("Iface", "gNMI"). // should be Id() as for builder
					Msg("Received gNMI response: " + prototext.Format(rsp.Response))

				if conf.GetBool("gnmi.publish") {
					if err := gh.KafkaHandler.PublishGNMINotificationToKafka(target, rsp.Response.String()); err != nil {
						errChan <- fmt.Errorf("failed to publish to Kafka: %v", err)
						return
					}
				}
			case tgErr := <-subErrChan:
				if tgErr.SubscriptionName == subscriptionID {
					log.Error().
						Str("Iface", "gNMI"). // should be Id() as for builder
						Msg("subscription " + tgErr.SubscriptionName + " stopped: " + tgErr.Err.Error())

					errChan <- fmt.Errorf("subscription %q stopped: %v", tgErr.SubscriptionName, tgErr.Err)
					return
				}
			case <-subCtx.Done():
				log.Info().
					Str("Iface", "gNMI"). // should be Id() as for builder
					Msg("subscription done: " + subscriptionID)

				errChan <- subCtx.Err()
				return
			}
		}
	}()

	// Wait for any errors
	//return <-errChan
	return subscriptionID, nil
}

// Unsubscribe stops a gNMI subscription.
func (gh *GNMIHandler) Unsubscribe(target string, subscriptionID string) error {
	if _, ok := nodeSubCtxCancelFuncs[target]; !ok {
		return fmt.Errorf("no subscriptions found for target %s", target)
	} else {
		// Stop all subscriptions for the target
		if subscriptionID == "" || subscriptionID == "*" {
			for subscriptionID, cancel := range nodeSubCtxCancelFuncs[target] {
				if cancel != nil {
					cancel()
					delete(nodeSubCtxCancelFuncs[target], subscriptionID)
				}
			}
			return nil
		} else {
			if _, ok := nodeSubCtxCancelFuncs[target][subscriptionID]; !ok {
				return fmt.Errorf("no subscription found for target %s and subscriptionID %s", target, subscriptionID)
			} else if nodeSubCtxCancelFuncs[target][subscriptionID] == nil {
				return fmt.Errorf("subscription %s already stopped", subscriptionID)
			}
			nodeSubCtxCancelFuncs[target][subscriptionID]()
			delete(nodeSubCtxCancelFuncs[target], subscriptionID)
			return nil
		}
	}
}
