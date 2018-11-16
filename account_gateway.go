package main

type AccountGateway interface {
	FindBy(userId string) (*Account, error)
	// Methods used in unit tests.
	Create(user User) (*Account, error)
	DeleteAll() error
}

func NewAccountGateway(database *Database) AccountGateway {
	return nil
}

type Account struct {
	Balance uint64
}
