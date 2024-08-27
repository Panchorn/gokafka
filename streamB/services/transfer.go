package services

import (
	"encoding/base64"
	"encoding/json"
	"encryption"
	"events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"time"
)

type EventHandler interface {
	Handle(topic string, payload []byte)
}

type eventHandler struct {
}

func NewEventHandler() EventHandler {
	return eventHandler{}
}

func (obj eventHandler) Handle(topic string, payload []byte) {
	log.Printf("handling topic %#v\n\n", topic)
	switch topic {
	case reflect.TypeOf(events.TransferExternalEvent{}).Name():
		createdEvent := &events.TransferExternalEvent{}
		err := json.Unmarshal(payload, createdEvent)
		if err != nil {
			log.Println(err)
			return
		}

		var callbackEvent events.Event

		secretToken, callbackEvent, ok := decodeSecretToken(createdEvent.RefID, createdEvent.SecretToken)
		if !ok {
			callbackEvent = events.TransferExternalFailedEvent{
				RefID:  createdEvent.RefID,
				Reason: "secret token is missing or invalid",
			}
		} else {
			log.Println("transfer is in progress with secretToken:", secretToken)
			time.Sleep(10 * time.Second)
			log.Println("transaction transferred")

			callbackEvent = events.TransferExternalCompletedEvent{
				RefID: createdEvent.RefID,
			}
		}

		producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
		if err != nil {
			panic(err)
		}
		defer producer.Close()

		producerHandler := NewEventProducer(producer)
		err = producerHandler.Produce(callbackEvent)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("message sent: %#v\n", callbackEvent)
	default:
		log.Println("topic unmatched", topic)
	}
}

func decodeSecretToken(refID string, secretToken string) (string, events.TransferExternalFailedEvent, bool) {
	callbackEvent := events.TransferExternalFailedEvent{
		RefID:  refID,
		Reason: "secret token is missing or invalid",
	}
	ciphertextDecoded, err := base64.StdEncoding.DecodeString(secretToken)
	if err != nil {
		return "", callbackEvent, false
	}
	plaintextDecrypted, err := encryption.Decrypt(ciphertextDecoded, encryption.Key())
	if err != nil {
		return "", callbackEvent, false
	}
	plaintext := string(plaintextDecrypted)
	if plaintext == "" {
		return "", callbackEvent, false
	}
	return plaintext, callbackEvent, true
}
