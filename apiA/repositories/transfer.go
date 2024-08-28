package repositories

import (
	"fmt"
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
	Save(transaction Transaction) error
	FindAll() (transactions []Transaction, err error)
	ExistsByRefID(refID string) (exists bool, err error)
	FindByRefID(refID string) (transaction Transaction, err error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	db.Table(tableNameTransactions).AutoMigrate(&Transaction{})
	return transactionRepository{db}
}

func (obj transactionRepository) Save(transaction Transaction) error {
	return obj.db.Table(tableNameTransactions).Save(transaction).Error
}

func (obj transactionRepository) FindAll() (transactions []Transaction, err error) {
	err = obj.db.Table(tableNameTransactions).Find(&transactions).Error
	return transactions, err
}

func (obj transactionRepository) ExistsByRefID(refID string) (exists bool, err error) {
	err = obj.db.Table(tableNameTransactions).Select("count(*) > 0").Where("ref_id=?", refID).Find(&exists).Error
	return exists, err
}

func (obj transactionRepository) FindByRefID(refID string) (transaction Transaction, err error) {
	err = obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).First(&transaction).Error
	return transaction, err
}

func (transaction Transaction) ToString() string {
	return fmt.Sprintf("%#v", transaction)
}
