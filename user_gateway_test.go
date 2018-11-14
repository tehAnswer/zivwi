package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

func TestFindBy(t *testing.T) {
	// Setting data and services.
	gateway := main.NewUserGateway()
	_, err1 := gateway.Create(main.User{
		FirstName: "Benito",
		LastName:  "Mussó",
		Email:     "benito@rome.it",
		Password:  "cia0p0rc0di0",
	})
	assert.NoError(t, err1)
	_, err2 := gateway.Create(main.User{
		FirstName: "Francisco",
		LastName:  "Paco",
		Email:     "paco@pp.es",
		Password:  "soydegal1c1a",
	})
	assert.NoError(t, err2)
	defer gateway.DeleteAll()

	returned, err := gateway.FindBy("benito@rome.it", "cia0p0rc0di0")
	if assert.NoError(t, err) {
		assert.Equal(t, "Benito", returned.FirstName)
		assert.Equal(t, "Mussó", returned.LastName)
	}

	returned, err = gateway.FindBy("paco@pp.es", "soydegal1c1a")
	if assert.NoError(t, err) {
		assert.Equal(t, "Francisco", returned.FirstName)
		assert.Equal(t, "Paco", returned.LastName)
	}

	returned, err = gateway.FindBy("adolfo@reich.de", "m3hrLi3b3")
	if assert.Error(t, err) {
		assert.Nil(t, returned)
	}
}
