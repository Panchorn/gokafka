package events

import "reflect"

var Topics = []string{
	reflect.TypeOf(TransferCreateEvent{}).Name(),
	reflect.TypeOf(TransferExternalEvent{}).Name(),
	reflect.TypeOf(TransferExternalCompletedEvent{}).Name(),
	reflect.TypeOf(TransferExternalFailedEvent{}).Name(),
	//reflect.TypeOf(OpenAccountEvent{}).Name(),
	//reflect.TypeOf(DepositFundEvent{}).Name(),
	//reflect.TypeOf(WithdrawFundEvent{}).Name(),
	//reflect.TypeOf(CloseAccountEvent{}).Name(),
}

type Event interface {
}

//type TransferEvent struct {
//	TransactionID     string
//	RefID  string
//	FromID string
//	ToID   string
//	Amount float64
//}

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

//type OpenAccountEvent struct {
//	TransactionID             string
//	AccountHolder  string
//	AccountType    int
//	OpeningBalance float64
//}
//
//type DepositFundEvent struct {
//	TransactionID     string
//	Amount float64
//}
//
//type WithdrawFundEvent struct {
//	TransactionID     string
//	Amount float64
//}
//
//type CloseAccountEvent struct {
//	TransactionID string
//}
