package services

import (
	"encoding/json"
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

		log.Println("transfer is in progress")
		time.Sleep(3 * time.Second)
		log.Println("transaction transferred")

		completedEvent := events.TransferExternalCompletedEvent{
			RefID: createdEvent.RefID,
		}

		producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
		if err != nil {
			panic(err)
		}
		defer producer.Close()

		producerHandler := NewEventProducer(producer)
		err = producerHandler.Produce(completedEvent)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("message sent: %#v\n", completedEvent)
	default:
		log.Println("topic unmatched", topic)
	}
}
