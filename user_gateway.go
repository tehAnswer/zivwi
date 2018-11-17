package main

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserGateway interface {
	FindBy(email string, password string) (*User, error)
	// Methods used in unit tests.
	Create(user User) (*User, error)
	DeleteAll() error
}

// Represents any person using the application.
type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UserGatewayImpl struct {
	Database *Database
}

func NewUserGateway(database *Database) UserGateway {
	return &UserGatewayImpl{
		Database: database,
	}
}

func (gtw *UserGatewayImpl) FindBy(email string, password string) (*User, error) {

	rows, dbError := gtw.Database.Connection.Query(`
    SELECT
       id, first_name, last_name, email, password
     FROM
       users
     WHERE
       email = $1`, email)

	if dbError != nil {
		return nil, dbError
	}

	var user User
	var scanErr error
	for rows.Next() {
		scanErr = rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
		)
	}

	// If there is no match, return nil and custom error.
	if user == (User{}) && scanErr == nil {
		return nil, fmt.Errorf("User not found")
	}

	pwdErr := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(password))
	if pwdErr != nil {
		return nil, fmt.Errorf("Incorrect email/password combination.")
	}

	return &user, scanErr
}

func (gtw *UserGatewayImpl) Create(user User) (*User, error) {
	query := `
    INSERT INTO users
      (id, first_name, last_name, email, password, created_at)
    VALUES
      ($1, $2, $3, $4, $5, $6)`
	saltedPassword := gtw.hashAndSalt(user.Password)
	userId := uuid.NewV4().String()
	_, dbError := gtw.Database.Connection.Query(query,
		userId,
		user.FirstName,
		user.LastName,
		user.Email,
		saltedPassword,
		time.Now(),
	)

	if dbError != nil {
		return nil, dbError
	}

	user.Id = userId
	return &user, nil
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
