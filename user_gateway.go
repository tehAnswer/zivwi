package main

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserGateway interface {
	FindBy(id string) (*User, error)
	// Methods used in unit tests.
	Create(firstName string, lastName string, password string) (*User, error)
	DeleteAll() error
}

// Represents any person using the application.
type User struct {
	Id        string
	FirstName string
	LastName  string
	Password  string
}

type UserGatewayImpl struct {
	Database *Database
}

func NewUserGateway() UserGateway {
	return &UserGatewayImpl{
		Database: NewDatabase(),
	}
}

func (gtw *UserGatewayImpl) FindBy(id string) (*User, error) {
	rows, dbError := gtw.Database.Connection.Query(`SELECT
     id, first_name, last_name FROM users where id = $1`, id)
	if dbError != nil {
		return nil, dbError
	}

	var user User
	var scanErr error
	for rows.Next() {
		scanErr = rows.Scan(&user.Id, &user.FirstName, &user.LastName)
	}
	// If there is no match, return nil and custom error.
	if user == (User{}) && scanErr == nil {
		return nil, fmt.Errorf("User not found")
	}

	return &user, scanErr
}

func (gtw *UserGatewayImpl) Create(firstName string, lastName string, password string) (*User, error) {
	query := `
    INSERT INTO users
      (id, first_name, last_name, password, created_at)
    VALUES
      ($1, $2, $3, $4, $5)`
	uuid := uuid.NewV4().String()
	saltedPassword := gtw.hashAndSalt(password)
	_, dbError := gtw.Database.Connection.Query(query,
		uuid,
		firstName,
		lastName,
		saltedPassword,
		time.Now(),
	)

	if dbError != nil {
		return nil, dbError
	}

	return &User{Id: uuid, FirstName: firstName, LastName: lastName}, nil
}

func (gtw *UserGatewayImpl) hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err.Error())
	}

	return string(hash)
}

func (gtw *UserGatewayImpl) DeleteAll() error {
	query := "DELETE FROM users"
	_, dbError := gtw.Database.Connection.Query(query)
	return dbError
}
