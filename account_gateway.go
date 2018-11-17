package main

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type AccountGateway interface {
	Where(accountIds []string) ([]*Account, error)
	FindBy(accountId string) (*Account, error)
	// Methods used in unit tests.
	Create(account Account) (*Account, error)
	DeleteAll() error
}

func NewAccountGateway(database *Database) AccountGateway {
	return &AccountGatewayImpl{
		Database: database,
	}
}

type AccountGatewayImpl struct {
	Database *Database
}

type Account struct {
	Id      string
	Balance uint64
}

func (gtw *AccountGatewayImpl) Where(accountIds []string) ([]*Account, error) {
	rows, dbError := gtw.Database.Connection.Query(`
      SELECT
         id, balance
       FROM
         accounts
       WHERE
         id = ANY($1)`, pq.Array(accountIds))
	defer rows.Close()
	if dbError != nil {
		return nil, dbError
	}

	var accounts []*Account
	var account Account
	var scanErr error

	for rows.Next() {
		scanErr = rows.Scan(&account.Id, &account.Balance)
		accounts = append(accounts, &Account{Id: account.Id, Balance: account.Balance})
	}

	return accounts, scanErr
}

func (gtw *AccountGatewayImpl) FindBy(accountId string) (*Account, error) {
	rows, dbError := gtw.Database.Connection.Query(`
      SELECT
         id, balance
       FROM
         accounts
       WHERE
         id = $1`, accountId)

	if dbError != nil {
		return nil, dbError
	}

	var account Account
	var scanErr error
	for rows.Next() {
		scanErr = rows.Scan(&account.Id, &account.Balance)
	}

	// If there is no match, return nil and custom error.
	if account == (Account{}) && scanErr == nil {
		return nil, fmt.Errorf("Account not found")
	}

	return &account, nil
}

func (gtw *AccountGatewayImpl) Create(account Account) (*Account, error) {
	query := `
    INSERT INTO accounts
      (id, balance, created_at)
    VALUES
      ($1, $2, $3)`
	accountId := uuid.NewV4().String()
	_, dbError := gtw.Database.Connection.Query(query,
		accountId,
		account.Balance,
		time.Now(),
	)

	if dbError != nil {
		return nil, dbError
	}

	account.Id = accountId
	return &account, nil
}

func (gtw *AccountGatewayImpl) DeleteAll() error {
	query := "DELETE FROM accounts"
	_, dbError := gtw.Database.Connection.Query(query)
	return dbError
}
