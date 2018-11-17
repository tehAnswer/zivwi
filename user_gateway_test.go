package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

func TestUserFindBy(t *testing.T) {
	// Setting data and services.
	gateway := main.NewUserGateway(main.NewDatabase())
	_, err1 := gateway.Create(main.User{
		FirstName:  "Benito",
		LastName:   "Mussó",
		Email:      "benito@rome.it",
		Password:   "cia0p0rc0di0",
		AccountIds: []string{"f44b7d9c-29c7-406b-99b2-159e036a28d6"},
	})
	assert.NoError(t, err1)
	_, err2 := gateway.Create(main.User{
		FirstName:  "Francisco",
		LastName:   "Paco",
		Email:      "paco@pp.es",
		Password:   "soydegal1c1a",
		AccountIds: []string{},
	})
	assert.NoError(t, err2)
	defer gateway.DeleteAll()

	returned, err := gateway.FindBy("benito@rome.it", "cia0p0rc0di0")
	if assert.NoError(t, err) {
		assert.Equal(t, "Benito", returned.FirstName)
		assert.Equal(t, "Mussó", returned.LastName)
		assert.Equal(t, "f44b7d9c-29c7-406b-99b2-159e036a28d6", returned.AccountIds[0])
	}

	returned, err = gateway.FindBy("paco@pp.es", "soydegal1c1a")
	if assert.NoError(t, err) {
		assert.Equal(t, "Francisco", returned.FirstName)
		assert.Equal(t, "Paco", returned.LastName)
		assert.Empty(t, returned.AccountIds)
	}

	returned, err = gateway.FindBy("adolfo@dhl.de", "m3hrLi3b3")
	if assert.Error(t, err) {
		assert.Nil(t, returned)
	}
}
