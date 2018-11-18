package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	app "github.com/tehAnswer/zivwi"
)

type FakeQueue struct{}

func (_ *FakeQueue) Publish(_ string, _ []byte) error {
	return nil
}

func TestValidTransfer(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	service := app.NewTransferService(
		accountGateway,
		app.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(app.Account{
		Balance: 10000,
	})

	transfer, err := service.Perform(account1.Id, "some_account_id", 10000, "xd")
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(10000), transfer.Amount)
		assert.Equal(t, account1.Id, transfer.FromAccountId)
		assert.Equal(t, "some_account_id", transfer.ToAccountId)
		assert.Equal(t, "xd", transfer.Message)
	}
}

func TestInvalidTransferDueToInsufficentFunds(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	service := app.NewTransferService(
		accountGateway,
		app.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(app.Account{
		Balance: 10000,
	})

	transfer, err := service.Perform(account1.Id, "some_account_id", 10001, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}

func TestValidTransferAndUpdateBalance(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	transferGateway := app.NewTransferGateway(database)

	defer accountGateway.DeleteAll()
	defer transferGateway.DeleteAll()
	service := app.NewTransferService(
		accountGateway,
		transferGateway,
		&FakeQueue{})

	account1, _ := accountGateway.Create(app.Account{
		Balance: 1000,
	})

	account2, _ := accountGateway.Create(app.Account{
		Balance: 1000,
	})

	transferGateway.Create(app.Transfer{
		FromAccountId: "",
		ToAccountId:   account1.Id,
		Amount:        1000,
		Status:        "completed",
	})

	transferGateway.Create(app.Transfer{
		FromAccountId: "",
		ToAccountId:   account2.Id,
		Amount:        1000,
		Status:        "completed",
	})

	transfer, err := service.Perform(account1.Id, account2.Id, 1, "xd")
	if assert.NoError(t, err) {
		assert.NotNil(t, transfer)
		resolveErr := service.Resolve(*transfer)
		assert.NoError(t, resolveErr)

		transfer, _ = transferGateway.FindBy(transfer.Id)
		assert.Equal(t, "completed", transfer.Status)

		account1, _ = accountGateway.FindBy(account1.Id)
		assert.Equal(t, uint64(999), account1.Balance)

		account2, _ = accountGateway.FindBy(account2.Id)
		assert.Equal(t, uint64(1001), account2.Balance)
	}
}

func TestInvalidTransferZeroCents(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	service := app.NewTransferService(
		accountGateway,
		app.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(app.Account{
		Balance: 10000,
	})

	transfer, err := service.Perform(account1.Id, account1.Id, 0, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}

func TestInvalidTransferNotSourceAccount(t *testing.T) {
	database := app.NewDatabase()
	accountGateway := app.NewAccountGateway(database)
	service := app.NewTransferService(
		accountGateway,
		app.NewTransferGateway(database),
		&FakeQueue{},
	)
	transfer, err := service.Perform("id_that_doesnt_exist", "some_id", 1000, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}
