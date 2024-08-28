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
	Handle(topic string, payload []byte)
}

type transferEventHandler struct {
	transferRepository repositories.TransactionRepository
	redis              TransactionService
}

func NewTransferEventHandler(transferRepository repositories.TransactionRepository, redis TransactionService) EventHandler {
	return transferEventHandler{transferRepository, redis}
}

func (obj transferEventHandler) Handle(topic string, payload []byte) {
	logs.Info("handling topic " + topic)
	switch topic {
	case reflect.TypeOf(events.TransferCreateEvent{}).Name():
		event := &events.TransferCreateEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			logs.Error(err)
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
		err = producerHandler.Produce(transferExternalEvent)
		if err != nil {
			logs.Error(err)
			return
		}
		logs.Info("message sent: " + transferExternalEvent.ToString())
	case reflect.TypeOf(events.TransferExternalCompletedEvent{}).Name():
		event := &events.TransferExternalCompletedEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			logs.Error(err)
			return
		}
		data := map[string]interface{}{
			"status":       "COMPLETED",
			"remark":       "transaction completed",
			"updated_date": time.Now(),
		}
		err = obj.transferRepository.PatchTransaction(event.RefID, data)
		if err != nil {
			logs.Error(err)
			return
		}
		err = evictTransaction(obj, event.RefID)
		if err != nil {
			logs.Error(err)
			return
		}
		logs.Info("patched transaction to COMPLETED")
	case reflect.TypeOf(events.TransferExternalFailedEvent{}).Name():
		event := &events.TransferExternalFailedEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			logs.Error(err)
			return
		}
		data := map[string]interface{}{
			"status":       "FAILED",
			"remark":       event.Reason,
			"updated_date": time.Now(),
		}
		err = obj.transferRepository.PatchTransaction(event.RefID, data)
		if err != nil {
			logs.Error(err)
			return
		}
		err = evictTransaction(obj, event.RefID)
		if err != nil {
			logs.Error(err)
			return
		}
		logs.Info("patched transaction to FAILED")
	default:
		logs.Info("topic unmatched: " + topic)
	}
}

func evictTransaction(obj transferEventHandler, refID string) error {
	err := obj.redis.EvictTransaction(refID)
	if err != nil {
		logs.Error(err)
		return err
	}
	logs.Info("redis transaction evicted")
	return nil
}
