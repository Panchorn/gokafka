package services

import (
	"encoding/json"
	"events"
	"github.com/IBM/sarama"
	"logs"
	"reflect"
)

type EventProducer interface {
	Produce(requestID string, event events.Event) error
}

type eventProducer struct {
	producer sarama.SyncProducer
}

func NewEventProducer(producer sarama.SyncProducer) EventProducer {
	return eventProducer{producer}
}

func (obj eventProducer) Produce(requestID string, event events.Event) error {
	topic := reflect.TypeOf(event).Name()
	logs.Debug(requestID, "producing message in topic "+topic)

	value, err := json.Marshal(event)
	if err != nil {
		logs.Error(requestID, err)
		return err
	}

	var eventHeaders []sarama.RecordHeader
	eventHeaders = append(eventHeaders, sarama.RecordHeader{
		Key:   sarama.ByteEncoder(logs.RequestID),
		Value: sarama.ByteEncoder(requestID),
	})

	msg := sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.ByteEncoder(requestID),
		Value:   sarama.ByteEncoder(value),
		Headers: eventHeaders,
	}

	_, _, err = obj.producer.SendMessage(&msg)
	if err != nil {
		logs.Error(requestID, err)
		return err
	}
	return nil
}
