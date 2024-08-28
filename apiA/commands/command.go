package commands

type TransferCommand struct {
	RefID       string  `json:"refID"`
	FromID      string  `json:"fromID"`
	ToID        string  `json:"toID"`
	Amount      float64 `json:"amount"`
	SecretToken string  `json:"secretToken"`
}
