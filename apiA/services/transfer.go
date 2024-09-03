package services

import (
	"apiA/commands"
	"apiA/repositories"
	"encoding/base64"
	"encryption"
	"errors"
	"events"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"logs"
	"time"
)

var poll = time.Millisecond * 100

type TransferService interface {
	Transfer(ctx echo.Context, command commands.TransferCommand) error
	TransferTransactions(ctx echo.Context) ([]repositories.Transaction, error)
}

type transferService struct {
	eventProducer EventProducer
	repository    repositories.TransactionRepository
	redis         TransactionService
}

func NewTransferService(eventProducer EventProducer, repository repositories.TransactionRepository, redis TransactionService) TransferService {
	return transferService{eventProducer, repository, redis}
}

func (obj transferService) Transfer(ctx echo.Context, command commands.TransferCommand) error {
	requestID := ctx.Get(logs.RequestID).(string)
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

	cipher, err := encryption.Encrypt([]byte(command.SecretToken), encryption.Key())
	if err != nil {
		return errors.New("error while encrypting secret token")
	}

	secretToken := base64.StdEncoding.EncodeToString(cipher)

	event := events.TransferCreateEvent{
		TransactionID: uuid.NewString(),
		RefID:         command.RefID,
		FromID:        command.FromID,
		ToID:          command.ToID,
		Amount:        command.Amount,
		SecretToken:   secretToken,
	}

	transaction := repositories.Transaction{
		ID:          event.TransactionID,
		RefID:       event.RefID,
		Status:      "AWAITING",
		FromID:      event.FromID,
		ToID:        event.ToID,
		Amount:      event.Amount,
		SecretToken: event.SecretToken,
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
	}
	err = obj.repository.Save(transaction)
	if err != nil {
		logs.Error(requestID, err)
		return errors.New("failed to save transfer transaction")
	}
	logs.Info(requestID, "saved transaction")

	logs.Info(requestID, "event"+event.ToString())
	eventHeaders := []events.EventHeader{
		{logs.RequestID, ctx.Get(logs.RequestID).(string)},
	}
	err = obj.eventProducer.Produce(ctx, event, eventHeaders)
	if err != nil {
		logs.Error(requestID, err)
		return errors.New("failed to produce event")
	}

	ticker := time.NewTicker(poll)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			transaction, err := obj.redis.GetTransaction(requestID, transaction.RefID)
			if err != nil {
				logs.Error(requestID, err)
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

func (obj transferService) TransferTransactions(ctx echo.Context) ([]repositories.Transaction, error) {
	requestID := ctx.Get(logs.RequestID).(string)
	transactions, err := obj.repository.FindAll()
	if err != nil {
		logs.Error(requestID, err)
		return nil, errors.New("failed to find transactions")
	}
	return transactions, nil
}
