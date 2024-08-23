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
	Save(transaction Transaction) error
	PatchStatus(refID string, status string) error
	FindAll() (transactions []Transaction, err error)
	FindByID(id string) (transaction Transaction, err error)
	//FindByRefID(refID string) (transaction Transaction, err error)
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

func (obj transactionRepository) PatchStatus(refID string, status string) error {
	return obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).Update("status", status).Error
}

func (obj transactionRepository) FindAll() (transactions []Transaction, err error) {
	err = obj.db.Table(tableNameTransactions).Find(&transactions).Error
	return transactions, err
}

func (obj transactionRepository) FindByID(id string) (transaction Transaction, err error) {
	err = obj.db.Table(tableNameTransactions).Where("id=?", id).First(&transaction).Error
	return transaction, err
}

//func (obj transactionRepository) FindByRefID(refID string) (transaction Transaction, err error) {
//	err = obj.db.Table(tableNameTransactions).Where("ref_id=?", refID).First(&transaction).Error
//	return transaction, err
//}

//type BankAccount struct {
//	TransactionID            string
//	AccountHolder string
//	AccountType   int
//	Balance       float64
//}

//type AccountRepository interface {
//Save(bankAccount BankAccount) error
//Delete(id string) error
//FindAll() (bankAccounts []BankAccount, err error)
//FindByID(id string) (bankAccount BankAccount, err error)
//}
