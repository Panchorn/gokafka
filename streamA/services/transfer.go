package services

import (
	"encoding/json"
	"events"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"streamA/repositories"
)

type EventHandler interface {
	Handle(topic string, payload []byte)
}

type transferEventHandler struct {
	transferRepository repositories.TransactionRepository
}

func NewTransferEventHandler(transferRepository repositories.TransactionRepository) EventHandler {
	return transferEventHandler{transferRepository}
}

func (obj transferEventHandler) Handle(topic string, payload []byte) {
	log.Printf("handling topic %#v\n", topic)
	switch topic {
	case reflect.TypeOf(events.TransferCreateEvent{}).Name():
		event := &events.TransferCreateEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			log.Println(err)
			return
		}
		id, err := uuid.NewUUID()
		if err != nil {
			log.Println(err)
			return
		}
		transaction := repositories.Transaction{
			ID:     id.String(),
			RefID:  event.RefID,
			Status: "AWAITING",
			FromID: event.FromID,
			ToID:   event.ToID,
			Amount: event.Amount,
		}
		err = obj.transferRepository.Save(transaction)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("saved transaction")

		transferExternalEvent := events.TransferExternalEvent{
			RefID:  event.RefID,
			FromID: event.FromID,
			ToID:   event.ToID,
			Amount: event.Amount,
		}

		producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
		if err != nil {
			panic(err)
		}
		defer producer.Close()

		producerHandler := NewEventProducer(producer)
		err = producerHandler.Produce(transferExternalEvent)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("message sent: %#v\n", transferExternalEvent)
	case reflect.TypeOf(events.TransferExternalCompletedEvent{}).Name():
		event := &events.TransferExternalCompletedEvent{}
		err := json.Unmarshal(payload, event)
		if err != nil {
			log.Println(err)
			return
		}
		err = obj.transferRepository.PatchStatus(event.RefID, "COMPLETED")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("patched transaction to COMPLETED")
	default:
		log.Println("topic unmatched", topic)
	}
}
