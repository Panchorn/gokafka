package services

import (
	"encoding/json"
	"events"
	"github.com/IBM/sarama"
	"log"
	"reflect"
)

type EventProducer interface {
	Produce(event events.Event) error
}

type eventProducer struct {
	producer sarama.SyncProducer
}

func NewEventProducer(producer sarama.SyncProducer) EventProducer {
	return eventProducer{producer}
}

func (obj eventProducer) Produce(event events.Event) error {
	topic := reflect.TypeOf(event).Name()
	log.Printf("producing message in topic %s", topic)

	value, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = obj.producer.SendMessage(&msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
