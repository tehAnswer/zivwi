package main

// The AuthorizationService is responsible for the user login logic. If a login
// is successful, the service returns a JWT and by the contrary, an error.

type AuthorizeService interface {
	login(email string, password string) (string, error)
}

type AuthorizeServiceImpl struct {
	Users      *UserGateway
	SigningKey []byte
}

func NewAuthorizeService() AuthorizeService {
	return nil
}
