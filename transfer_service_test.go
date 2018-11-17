package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

type FakeQueue struct{}

func (_ *FakeQueue) Publish(_ string, _ []byte) error {
	return nil
}

func TestValidTransfer(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	service := main.NewTransferService(
		accountGateway,
		main.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(main.Account{
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
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	service := main.NewTransferService(
		accountGateway,
		main.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(main.Account{
		Balance: 10000,
	})

	transfer, err := service.Perform(account1.Id, "some_account_id", 10001, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}

func TestInvalidTransferSameAccount(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	service := main.NewTransferService(
		accountGateway,
		main.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(main.Account{
		Balance: 10000,
	})

	transfer, err := service.Perform(account1.Id, account1.Id, 1000, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}

func TestInvalidTransferZeroCents(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	service := main.NewTransferService(
		accountGateway,
		main.NewTransferGateway(database),
		&FakeQueue{},
	)

	account1, _ := accountGateway.Create(main.Account{
		Balance: 10000,
	})

	transfer, err := service.Perform(account1.Id, account1.Id, 0, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}

func TestInvalidTransferNotSourceAccount(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	service := main.NewTransferService(
		accountGateway,
		main.NewTransferGateway(database),
		&FakeQueue{},
	)
	transfer, err := service.Perform("id_that_doesnt_exist", "some_id", 1000, "xd")
	assert.Error(t, err)
	assert.Nil(t, transfer)
}
