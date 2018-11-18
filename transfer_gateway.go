package app

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransferGateway interface {
	UpdateBalancesFrom(transfer Transfer) error
	Create(transfer Transfer) (*Transfer, error)
	FindBy(transferId string) (*Transfer, error)
	Update(transfer Transfer) (*Transfer, error)
	DeleteAll() error
}

type TransferGatewayImpl struct {
	Database *Database
}

type Transfer struct {
	Id            string `json:"id"`
	FromAccountId string `json:"from_account_id"`
	ToAccountId   string `json:"to_account_id"`
	Message       string `json:"message"`
	Amount        uint64 `json:"amount"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
}

func NewTransferGateway(database *Database) TransferGateway {
	return &TransferGatewayImpl{
		Database: database,
	}
}

func (gtw *TransferGatewayImpl) Create(transfer Transfer) (*Transfer, error) {
	query := `
    INSERT INTO transfers
      (id, from_account_id, to_account_id, message, amount, status, created_at, updated_at)
    VALUES
      ($1, $2, $3, $4, $5, $6, $7, $8)`
	transferId := uuid.NewV4().String()
	timestamp := time.Now()
	_, dbError := gtw.Database.Connection.Query(query,
		transferId,
		transfer.FromAccountId,
		transfer.ToAccountId,
		transfer.Message,
		transfer.Amount,
		transfer.Status,
		timestamp,
		timestamp,
	)

	if dbError != nil {
		return nil, dbError
	}

	transfer.Id = transferId
	return &transfer, nil
}

func (gtw *TransferGatewayImpl) FindBy(transferId string) (*Transfer, error) {
	rows, dbError := gtw.Database.Connection.Query(`
      SELECT
         id, from_account_id, to_account_id, message, amount, status, error
       FROM
         transfers
       WHERE
         id = $1`, transferId)

	if dbError != nil {
		return nil, dbError
	}

	var transfer Transfer
	var scanErr error

	for rows.Next() {
		scanErr = rows.Scan(
			&transfer.Id,
			&transfer.FromAccountId,
			&transfer.ToAccountId,
			&transfer.Message,
			&transfer.Amount,
			&transfer.Status,
			&transfer.Error,
		)
	}

	// If there is no match, return nil and custom error.
	if transfer == (Transfer{}) && scanErr == nil {
		return nil, fmt.Errorf("Transfer not found")
	}

	return &transfer, nil
}

func (gtw *TransferGatewayImpl) Update(transfer Transfer) (*Transfer, error) {
	tx, txErr := gtw.Database.Connection.Begin()
	if txErr != nil {
		return nil, txErr
	}
	return gtw.update(tx, transfer)
}

func (gtw *TransferGatewayImpl) update(tx *sql.Tx, transfer Transfer) (*Transfer, error) {
	updatedAt := time.Now()

	rows, dbError := tx.Query(`
      UPDATE transfers
      SET status = $2, error = $3, updated_at = $4
      WHERE id = $1
			RETURNING id`, transfer.Id, transfer.Status, transfer.Error, updatedAt)

	if dbError != nil {
		return nil, dbError
	}

	updatesCheckErr := EnsureOneUpdate(rows)
	if updatesCheckErr != nil {
		return nil, updatesCheckErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, commitErr
	}

	return &transfer, nil
}

func (gtw *TransferGatewayImpl) DeleteAll() error {
	query := "DELETE FROM transfers"
	_, dbError := gtw.Database.Connection.Query(query)
	return dbError
}

func (gtw *TransferGatewayImpl) UpdateBalancesFrom(transfer Transfer) error {

	tx, txErr := gtw.Database.Connection.Begin()
	if txErr != nil {
		return txErr
	}

	fromAccountBalancingErr := UpdateBalanceOf(tx, transfer.FromAccountId, transfer.Id)
	if fromAccountBalancingErr != nil {
		return fromAccountBalancingErr
	}
	toAccountBalancingErr := UpdateBalanceOf(tx, transfer.ToAccountId, transfer.Id)
	if toAccountBalancingErr != nil {
		return toAccountBalancingErr
	}
	transfer.Status = "completed"

	gtw.update(tx, transfer)
	return nil
}

func UpdateBalanceOf(tx *sql.Tx, accountId string, transferId string) error {
	query := `
		UPDATE accounts
		SET balance = subquery.balance
		FROM (
			SELECT
				sum(
					CASE
					WHEN from_account_id = $1 THEN amount * -1
					WHEN to_account_id = $1 THEN amount
					END
				) AS balance
				FROM transfers
				WHERE from_account_id = $1
				OR to_account_id = $1
				AND status = 'completed'
				OR id = $2
		) AS subquery
		WHERE id = $1
		RETURNING id`

	rows, errorBalance := tx.Query(query, accountId, transferId)
	if errorBalance != nil {
		fmt.Println(errorBalance.Error())
		return errorBalance
	}

	oneUpdateError := EnsureOneUpdate(rows)
	if oneUpdateError != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to rebalance %v. Reason: %v",
			accountId,
			oneUpdateError.Error())
	}
	return nil
}

func EnsureOneUpdate(rows *sql.Rows) error {
	counter := 0
	for rows.Next() {
		counter = counter + 1
	}

	if counter != 1 {
		return fmt.Errorf("None or more than one record were updated.")
	}
	return nil
}
