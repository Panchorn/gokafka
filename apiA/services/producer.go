package services

import (
	"encoding/json"
	"events"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"logs"
	"reflect"
)

type EventProducer interface {
	Produce(ctx echo.Context, event events.Event, headers []events.EventHeader) error
}

type eventProducer struct {
	producer sarama.SyncProducer
}

func NewEventProducer(producer sarama.SyncProducer) EventProducer {
	return eventProducer{producer}
}

func (obj eventProducer) Produce(ctx echo.Context, event events.Event, headers []events.EventHeader) error {
	requestID := ctx.Get(logs.RequestID).(string)
	topic := reflect.TypeOf(event).Name()
	logs.Debug(requestID, "producing message in topic "+topic)

	value, err := json.Marshal(event)
	if err != nil {
		logs.Error(requestID, err)
		return err
	}

	var eventHeaders []sarama.RecordHeader
	for _, header := range headers {
		eventHeaders = append(eventHeaders, sarama.RecordHeader{
			Key:   sarama.ByteEncoder(header.Key),
			Value: sarama.ByteEncoder(header.Value),
		})
	}

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
