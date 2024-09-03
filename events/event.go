package events

import (
	"fmt"
	"reflect"
)

var Topics = []string{
	reflect.TypeOf(TransferCreateEvent{}).Name(),
	reflect.TypeOf(TransferExternalEvent{}).Name(),
	reflect.TypeOf(TransferExternalCompletedEvent{}).Name(),
	reflect.TypeOf(TransferExternalFailedEvent{}).Name(),
}

type Event interface {
	ToString() string
}

type EventHeader struct {
	Key   string
	Value string
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

func (event TransferCreateEvent) ToString() string {
	return fmt.Sprintf("%#v", event)
}

func (event TransferExternalEvent) ToString() string {
	return fmt.Sprintf("%#v", event)
}

func (event TransferExternalCompletedEvent) ToString() string {
	return fmt.Sprintf("%#v", event)
}

func (event TransferExternalFailedEvent) ToString() string {
	return fmt.Sprintf("%#v", event)
}
