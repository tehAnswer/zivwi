package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

func TestFindBy(t *testing.T) {
	// Setting data and services.
	gateway := main.NewUserGateway()
	user1, err1 := gateway.Create("Benito", "Mussó", "ciaoporcodi0")
	assert.NoError(t, err1)
	user2, err2 := gateway.Create("Francisco", "Paco", "soydegalicia1")
	assert.NoError(t, err2)
	defer gateway.DeleteAll()

	returned, err := gateway.FindBy(user1.Id)
	if assert.NoError(t, err) {
		assert.Equal(t, "Benito", returned.FirstName)
		assert.Equal(t, "Mussó", user1.LastName)
	}

	returned, err = gateway.FindBy(user2.Id)
	if assert.NoError(t, err) {
		assert.Equal(t, "Francisco", returned.FirstName)
		assert.Equal(t, "Paco", returned.LastName)
	}

	returned, err = gateway.FindBy("not_exist")
	if assert.Error(t, err) {
		assert.Nil(t, returned)
	}
}
