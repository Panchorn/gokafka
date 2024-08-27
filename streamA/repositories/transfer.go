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
	Remark      string
	Amount      float64
	SecretToken string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type TransactionRepository interface {
	PatchTransaction(refID string, data interface{}) error
	FindByRefID(refID string) (transaction Transaction, err error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	db.Table(tableNameTransactions).AutoMigrate(&Transaction{})
	return transactionRepository{db}
}

func (obj transactionRepository) PatchTransaction(refID string, data interface{}) error {
	return obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).Updates(data).Error
}

func (obj transactionRepository) FindByRefID(refID string) (transaction Transaction, err error) {
	err = obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).First(&transaction).Error
	return transaction, err
}
