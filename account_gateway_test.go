package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

func TestAccountFindBy(t *testing.T) {
	// Setting data and services.
	gateway := main.NewAccountGateway(main.NewDatabase())
	account, err := gateway.Create(main.Account{Balance: uint64(10000)})
	assert.NoError(t, err)
	defer gateway.DeleteAll()

	returned, err := gateway.FindBy(account.Id)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(10000), returned.Balance)
	}
}
func TestAccountWhere(t *testing.T) {
	// Setting data and services.
	gateway := main.NewAccountGateway(main.NewDatabase())
	account1, err1 := gateway.Create(main.Account{Balance: uint64(10000)})
	assert.NoError(t, err1)
	account2, err2 := gateway.Create(main.Account{Balance: uint64(20000)})
	assert.NoError(t, err2)
	defer gateway.DeleteAll()

	ids := []string{account1.Id, account2.Id}
	returned, err := gateway.Where(ids)
	if assert.NoError(t, err) {
		assert.Len(t, returned, 2)
		assert.Equal(t, uint64(10000), returned[0].Balance)
		assert.Equal(t, uint64(20000), returned[1].Balance)
	}
}
