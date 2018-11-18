package app

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// The AuthorizationService is responsible for the user login logic. If a login
// is successful, the service returns a JWT and by the contrary, an error.

type AuthorizeService interface {
	Login(email string, password string) (string, error)
}

type AuthorizeServiceImpl struct {
	Users      UserGateway
	SigningKey []byte
}

func NewAuthorizeService(userGateway UserGateway) AuthorizeService {
	return &AuthorizeServiceImpl{
		Users:      userGateway,
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
}

func (service *AuthorizeServiceImpl) Login(email string, password string) (string, error) {
	user, loginErr := service.Users.FindBy(email, password)

	if loginErr == nil {
		return service.buildJwtFor(user)
	}
	return "", loginErr
}

func (service *AuthorizeServiceImpl) buildJwtFor(user *User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.Id
	claims["name"] = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	claims["role"] = "user"
	claims["account_ids"] = user.AccountIds
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Generate encoded token and send it as response.
	return token.SignedString(service.SigningKey)
}
