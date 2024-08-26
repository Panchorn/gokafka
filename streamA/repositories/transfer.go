package repositories

import (
	"gorm.io/gorm"
	"time"
)

var tableNameTransactions = "transactions"

type Transaction struct {
	ID          string
	RefID       string
	FromID      string
	ToID        string
	Status      string
	Amount      float64
	CreatedDate time.Time
	UpdatedDate time.Time
}

type TransactionRepository interface {
	PatchStatus(refID string, status string) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	db.Table(tableNameTransactions).AutoMigrate(&Transaction{})
	return transactionRepository{db}
}

func (obj transactionRepository) PatchStatus(refID string, status string) error {
	return obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).Update("status", status).Update("updated_date", time.Now()).Error
}
