package main

type UserGateway interface {
	findUserBy(id string) (*User, error)
}

// Represents any person using the application.
type User struct {
	Id        string
	FirstName string
	LastName  string
}
