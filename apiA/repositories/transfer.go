package repositories

import "gorm.io/gorm"

var tableNameTransactions = "transactions"

type Transaction struct {
	ID     string
	RefID  string
	FromID string
	ToID   string
	Status string
	Amount float64
}

type TransactionRepository interface {
	FindByRefID(refID string) (transaction Transaction, err error)
}

type transactionRepository struct {
	db *gorm.DB
}

func (obj transactionRepository) FindByRefID(refID string) (transaction Transaction, err error) {
	err = obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).First(&transaction).Error
	return transaction, err
}
