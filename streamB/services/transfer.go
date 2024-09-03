package services

import (
	"encoding/base64"
	"encoding/json"
	"encryption"
	"events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"logs"
	"reflect"
)

type EventHandler interface {
	Handle(topic string, key []byte, payload []byte, headers []*sarama.RecordHeader)
}

type eventHandler struct {
}

func NewEventHandler() EventHandler {
	return eventHandler{}
}

func (obj eventHandler) Handle(topic string, key []byte, payload []byte, headers []*sarama.RecordHeader) {
	requestID := string(key)
	logs.Info(requestID, "handling topic "+topic)
	switch topic {
	case reflect.TypeOf(events.TransferExternalEvent{}).Name():
		createdEvent := &events.TransferExternalEvent{}
		err := json.Unmarshal(payload, createdEvent)
		if err != nil {
			logs.Error(requestID, err)
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
			logs.Info(requestID, "transfer is in progress with secretToken "+secretToken)
			//time.Sleep(500 * time.Millisecond)
			logs.Info(requestID, "transaction transferred")

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
		err = producerHandler.Produce(requestID, callbackEvent)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		logs.Info(requestID, "message sent: "+callbackEvent.ToString())
	default:
		logs.Info(requestID, "topic unmatched "+topic)
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
