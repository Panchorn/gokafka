package events

import "reflect"

var Topics = []string{
	reflect.TypeOf(TransferCreateEvent{}).Name(),
	reflect.TypeOf(TransferExternalEvent{}).Name(),
	reflect.TypeOf(TransferExternalCompletedEvent{}).Name(),
	reflect.TypeOf(TransferExternalFailedEvent{}).Name(),
}

type Event interface {
}

type TransferCreateEvent struct {
	TransactionID string
	RefID         string
	FromID        string
	ToID          string
	Amount        float64
	SecretToken   string
}

type TransferExternalEvent struct {
	RefID       string
	FromID      string
	ToID        string
	Amount      float64
	SecretToken string
}

type TransferExternalCompletedEvent struct {
	RefID string
}

type TransferExternalFailedEvent struct {
	RefID  string
	Reason string
}
