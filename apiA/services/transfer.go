package services

import (
	"apiA/commands"
	"apiA/repositories"
	"errors"
	"events"
	"github.com/google/uuid"
	"log"
	"time"
)

var poll = time.Millisecond * 500

type TransferService interface {
	Transfer(command commands.TransferCommand) error
	TransferTransactions() ([]repositories.Transaction, error)
}

type transferService struct {
	eventProducer EventProducer
	repository    repositories.TransactionRepository
}

func NewTransferService(eventProducer EventProducer, repository repositories.TransactionRepository) TransferService {
	return transferService{eventProducer, repository}
}

func (obj transferService) Transfer(command commands.TransferCommand) error {
	if command.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	isDuplicateTransaction, err := obj.repository.ExistsByRefID(command.RefID)
	if err != nil {
		return errors.New("error while finding duplicate transaction")
	}
	if isDuplicateTransaction {
		return errors.New("duplicate transaction")
	}

	event := events.TransferCreateEvent{
		TransactionID: uuid.NewString(),
		RefID:         command.RefID,
		FromID:        command.FromID,
		ToID:          command.ToID,
		Amount:        command.Amount,
		SecretToken:   command.SecretToken,
	}

	transaction := repositories.Transaction{
		ID:          event.TransactionID,
		RefID:       event.RefID,
		Status:      "AWAITING",
		FromID:      event.FromID,
		ToID:        event.ToID,
		Amount:      event.Amount,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}
	err = obj.repository.Save(transaction)
	if err != nil {
		log.Println(err)
		return errors.New("failed to save transfer transaction")
	}
	log.Println("saved transaction")

	log.Printf("%#v", event)
	err = obj.eventProducer.Produce(event)
	if err != nil {
		log.Println(err)
		return errors.New("failed to produce event")
	}

	ticker := time.NewTicker(poll)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			transaction, err := obj.repository.FindByID(transaction.ID)
			if err != nil {
				log.Println(err)
				return errors.New("failed to fetch transaction")
			}
			if transaction.Status == "COMPLETED" {
				return nil
			}
			if transaction.Status == "FAILED" {
				return errors.New("transfer failed: " + transaction.Remark)
			}
		}
	}
}

func (obj transferService) TransferTransactions() ([]repositories.Transaction, error) {
	transactions, err := obj.repository.FindAll()
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to find transactions")
	}
	return transactions, nil
}
