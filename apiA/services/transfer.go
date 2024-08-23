package services

import (
	"apiA/commands"
	"errors"
	"events"
	"github.com/google/uuid"
	"log"
)

type TransferService interface {
	Transfer(command commands.TransferCommand) error
}

type transferService struct {
	eventProducer EventProducer
}

func NewTransferService(eventProducer EventProducer) TransferService {
	return transferService{eventProducer}
}

func (obj transferService) Transfer(command commands.TransferCommand) error {
	if command.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	event := events.TransferCreateEvent{
		TransactionID: uuid.NewString(),
		RefID:         command.RefID,
		FromID:        command.FromID,
		ToID:          command.ToID,
		Amount:        command.Amount,
	}

	log.Printf("%#v", event)
	return obj.eventProducer.Produce(event)
}
