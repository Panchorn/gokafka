package services

import (
	"encoding/json"
	"events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"logs"
	"reflect"
	"streamA/repositories"
	"time"
)

type EventHandler interface {
	Handle(topic string, key []byte, payload []byte, headers []*sarama.RecordHeader)
}

type transferEventHandler struct {
	transferRepository repositories.TransactionRepository
	redis              TransactionService
}

func NewTransferEventHandler(transferRepository repositories.TransactionRepository, redis TransactionService) EventHandler {
	return transferEventHandler{transferRepository, redis}
}

func (obj transferEventHandler) Handle(topic string, key []byte, payload []byte, headers []*sarama.RecordHeader) {
	requestID := string(key)
	logs.Info(requestID, "handling topic "+topic)
	for _, header := range headers {
		logs.Info(requestID, "handling topic with header "+string(header.Key)+" "+string(header.Value))
	}

	switch topic {
	case reflect.TypeOf(events.TransferCreateEvent{}).Name():
		event := &events.TransferCreateEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			logs.Error(requestID, err)
			return
		}

		transferExternalEvent := events.TransferExternalEvent{
			RefID:       event.RefID,
			FromID:      event.FromID,
			ToID:        event.ToID,
			Amount:      event.Amount,
			SecretToken: event.SecretToken,
		}

		producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
		if err != nil {
			panic(err)
		}
		defer producer.Close()

		producerHandler := NewEventProducer(producer)
		err = producerHandler.Produce(requestID, transferExternalEvent)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		logs.Info(requestID, "message sent: "+transferExternalEvent.ToString())
	case reflect.TypeOf(events.TransferExternalCompletedEvent{}).Name():
		event := &events.TransferExternalCompletedEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		data := map[string]interface{}{
			"status":       "COMPLETED",
			"remark":       "transaction completed",
			"updated_date": time.Now(),
		}
		err = obj.transferRepository.PatchTransaction(event.RefID, data)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		err = evictTransaction(obj, requestID, event.RefID)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		logs.Info(requestID, "patched transaction to COMPLETED")
	case reflect.TypeOf(events.TransferExternalFailedEvent{}).Name():
		event := &events.TransferExternalFailedEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		data := map[string]interface{}{
			"status":       "FAILED",
			"remark":       event.Reason,
			"updated_date": time.Now(),
		}
		err = obj.transferRepository.PatchTransaction(event.RefID, data)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		err = evictTransaction(obj, requestID, event.RefID)
		if err != nil {
			logs.Error(requestID, err)
			return
		}
		logs.Info(requestID, "patched transaction to FAILED")
	default:
		logs.Info(requestID, "topic unmatched: "+topic)
	}
}

func evictTransaction(obj transferEventHandler, requestID string, refID string) error {
	err := obj.redis.EvictTransaction(refID)
	if err != nil {
		logs.Error(requestID, err)
		return err
	}
	logs.Info(requestID, "redis transaction evicted")
	return nil
}
