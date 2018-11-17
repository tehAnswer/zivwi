package main_test

import (
	"fmt"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

func TestCorrectLogin(t *testing.T) {
	service := main.NewAuthorizeService(
		main.NewUserGateway(main.NewDatabase()))
	authorizationService, _ := service.(*main.AuthorizeServiceImpl)
	user, _ := authorizationService.Users.Create(main.User{
		FirstName:  "Benito",
		LastName:   "Mussó",
		Email:      "benito@rome.it",
		Password:   "cia0p0rc0di0",
		AccountIds: []string{"b2cc7786-a4ce-45b0-9268-75aed2f14554"},
	})

	defer authorizationService.Users.DeleteAll()

	token, loginErr := service.Login("benito@rome.it", "cia0p0rc0di0")
	if assert.NoError(t, loginErr) {
		parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		assert.True(t, parsedToken.Valid)
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			assert.Equal(t, user.Id, claims["sub"])
			assert.Equal(t, "Benito Mussó", claims["name"])
			assert.Equal(t, "user", claims["role"])
			assert.NotNil(t, claims["exp"])
			ids, _ := claims["account_ids"].([]interface{})
			assert.Equal(t, "b2cc7786-a4ce-45b0-9268-75aed2f14554", ids[0].(string))
		}
	}

}

func TestIncorrectLogin(t *testing.T) {
	service := main.NewAuthorizeService(
		main.NewUserGateway(main.NewDatabase()))
	authorizationService, _ := service.(*main.AuthorizeServiceImpl)
	authorizationService.Users.Create(main.User{
		FirstName: "Benito",
		LastName:  "Mussó",
		Email:     "benito@rome.it",
		Password:  "cia0p0rc0di0",
	})

	defer authorizationService.Users.DeleteAll()

	token, loginErr := service.Login("benito@rome.it", "ciaobella")
	assert.Empty(t, token)
	assert.Error(t, loginErr)
}
