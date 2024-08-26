package commands

type TransferCommand struct {
	RefID       string
	FromID      string
	ToID        string
	Amount      float64
	SecretToken string
}
