package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

func TestTransferFindBy(t *testing.T) {
	// Setting data and services.
	database := main.NewDatabase()
	gateway := main.NewTransferGateway(database)
	transfer, err := gateway.Create(main.Transfer{
		FromAccountId: "3c489418-8d96-45c1-9773-321d21062964",
		ToAccountId:   "1e011113-d42c-4318-b9de-ee9cc064d4c0",
		Amount:        10000,
		Message:       "your stake daniel",
		Status:        "processing",
	})
	assert.NoError(t, err)
	defer gateway.DeleteAll()

	returned, err := gateway.FindBy(transfer.Id)

	if assert.NoError(t, err) {
		assert.Equal(t, "3c489418-8d96-45c1-9773-321d21062964", returned.FromAccountId)
		assert.Equal(t, "1e011113-d42c-4318-b9de-ee9cc064d4c0", returned.ToAccountId)
		assert.Equal(t, "processing", returned.Status)
		assert.Equal(t, "your stake daniel", returned.Message)
		assert.Equal(t, uint64(10000), returned.Amount)
	}
}

func TestTransferUpdate(t *testing.T) {
	// Setting data and services.
	database := main.NewDatabase()
	gateway := main.NewTransferGateway(database)
	transfer, err := gateway.Create(main.Transfer{
		FromAccountId: "3c489418-8d96-45c1-9773-321d21062964",
		ToAccountId:   "1e011113-d42c-4318-b9de-ee9cc064d4c0",
		Amount:        10000,
		Message:       "Beers for the party",
		Status:        "processing",
	})
	assert.NoError(t, err)
	defer gateway.DeleteAll()

	transfer.Status = "cancelled"
	transfer.Error = "not_enough_funds"

	_, updateErr := gateway.Update(*transfer)
	assert.NoError(t, updateErr)
	returned, _ := gateway.FindBy(transfer.Id)

	assert.Equal(t, "3c489418-8d96-45c1-9773-321d21062964", returned.FromAccountId)
	assert.Equal(t, "1e011113-d42c-4318-b9de-ee9cc064d4c0", returned.ToAccountId)
	assert.Equal(t, "cancelled", returned.Status)
	assert.Equal(t, "not_enough_funds", returned.Error)
	assert.Equal(t, "Beers for the party", returned.Message)
	assert.Equal(t, uint64(10000), returned.Amount)
}
