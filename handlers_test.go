package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

var (
	correctCredentials   = `{"email": "benito@rome.it", "password":"cia0p0rc0di0"}`
	incorrectCredentials = `{"email": "benito@rome.it", "password":"ciao"}`
)

func TestLoginThroughAPIAndTransfer(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	transferGateway := main.NewTransferGateway(database)
	userGateway := main.NewUserGateway(database)

	account1, _ := accountGateway.Create(main.Account{
		Balance: 10000,
	})

	account2, _ := accountGateway.Create(main.Account{
		Balance: 10000000,
	})
	user, _ := userGateway.Create(main.User{
		FirstName: "Benito",
		LastName:  "Mussó",
		Email:     "benito@rome.it",
		Password:  "cia0p0rc0di0",
		AccountIds: []string{
			account1.Id,
			account2.Id,
		},
	})

	defer transferGateway.DeleteAll()
	defer accountGateway.DeleteAll()
	defer userGateway.DeleteAll()

	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/authorize",
		strings.NewReader(correctCredentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions for login
	if assert.NoError(t, main.Authorize(c)) {
		assert.Equal(t, 200, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}

	var body struct {
		Token string `json:"token"`
	}

	json.Unmarshal(rec.Body.Bytes(), &body)

	payload, _ := json.Marshal(struct {
		FromAccountId string `json:"from_account_id"`
		ToAccountId   string `json:"to_account_id"`
		Message       string `json:"message"`
		Amount        uint64 `json:"amount"`
	}{user.AccountIds[1], user.AccountIds[0], "", uint64(1)})

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/cmd/transfer",
		strings.NewReader(string(payload)))

	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	parsedToken, _ := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	c2.Set("user", parsedToken)

	// Assertions for login
	if assert.NoError(t, main.TransferCmd(c2)) {
		assert.Equal(t, 201, rec2.Code)
		assert.Contains(t, rec2.Body.String(), "processing")
		var body2 struct {
			Id string `json:"id"`
		}
		json.Unmarshal(rec2.Body.Bytes(), &body2)
		transfer, findErr := transferGateway.FindBy(body2.Id)
		if assert.NoError(t, findErr) {
			assert.Equal(t, "processing", transfer.Status)
			assert.Equal(t, account2.Id, transfer.FromAccountId)
			assert.Equal(t, account1.Id, transfer.ToAccountId)
		}
	}
}

func TestIncorrectLoginThroughAPI(t *testing.T) {
	gateway := main.NewUserGateway(main.NewDatabase())
	_, err := gateway.Create(main.User{
		FirstName: "Benito",
		LastName:  "Mussó",
		Email:     "benito@rome.it",
		Password:  "cia0p0rc0di0",
	})
	defer gateway.DeleteAll()
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/authorize",
		strings.NewReader(incorrectCredentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, main.Authorize(c)) {
		assert.Equal(t, 401, rec.Code)
		assert.Contains(t, rec.Body.String(), "Incorrect email/password combination")
	}
}

func TestIncorrectLoginParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/authorize", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, main.Authorize(c)) {
		assert.Equal(t, 422, rec.Code)
		assert.Contains(t, rec.Body.String(), "Incorrect login params")
	}
}

func TestInvalidTransferDueToTenancy(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	transferGateway := main.NewTransferGateway(database)
	userGateway := main.NewUserGateway(database)

	account1, _ := accountGateway.Create(main.Account{
		Balance: 10000,
	})

	account2, _ := accountGateway.Create(main.Account{
		Balance: 10000000,
	})

	userGateway.Create(main.User{
		FirstName:  "Benito",
		LastName:   "Mussó",
		Email:      "benito@rome.it",
		Password:   "cia0p0rc0di0",
		AccountIds: []string{},
	})

	defer transferGateway.DeleteAll()
	defer accountGateway.DeleteAll()
	defer userGateway.DeleteAll()

	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/authorize",
		strings.NewReader(correctCredentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions for login
	if assert.NoError(t, main.Authorize(c)) {
		assert.Equal(t, 200, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}

	var body struct {
		Token string `json:"token"`
	}

	json.Unmarshal(rec.Body.Bytes(), &body)

	payload, _ := json.Marshal(struct {
		FromAccountId string `json:"from_account_id"`
		ToAccountId   string `json:"to_account_id"`
		Message       string `json:"message"`
		Amount        uint64 `json:"amount"`
	}{account1.Id, account2.Id, "", uint64(1)})

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/cmd/transfer",
		strings.NewReader(string(payload)))

	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	parsedToken, _ := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	c2.Set("user", parsedToken)

	// Assertions for login
	if assert.NoError(t, main.TransferCmd(c2)) {
		assert.Equal(t, 401, rec2.Code)
	}
}

func TestAccounts(t *testing.T) {
	database := main.NewDatabase()
	accountGateway := main.NewAccountGateway(database)
	transferGateway := main.NewTransferGateway(database)
	userGateway := main.NewUserGateway(database)

	account1, _ := accountGateway.Create(main.Account{
		Balance: 10000,
	})

	account2, _ := accountGateway.Create(main.Account{
		Balance: 10000000,
	})

	userGateway.Create(main.User{
		FirstName:  "Benito",
		LastName:   "Mussó",
		Email:      "benito@rome.it",
		Password:   "cia0p0rc0di0",
		AccountIds: []string{account1.Id, account2.Id},
	})

	defer transferGateway.DeleteAll()
	defer accountGateway.DeleteAll()
	defer userGateway.DeleteAll()

	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/authorize",
		strings.NewReader(correctCredentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions for login
	if assert.NoError(t, main.Authorize(c)) {
		assert.Equal(t, 200, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}

	var body struct {
		Token string `json:"token"`
	}

	json.Unmarshal(rec.Body.Bytes(), &body)

	req2 := httptest.NewRequest(
		http.MethodGet,
		"/q/accounts",
		strings.NewReader(""))

	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	parsedToken, _ := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	c2.Set("user", parsedToken)

	// Assertions for login
	if assert.NoError(t, main.AccountsQuery(c2)) {
		assert.Equal(t, 200, rec2.Code)

		var body2 struct {
			Accounts []*main.Account `json:"accounts"`
		}

		json.Unmarshal(rec2.Body.Bytes(), &body2)

		assert.Len(t, body2.Accounts, 2)
		assert.Equal(t, account1.Balance, body2.Accounts[0].Balance)
		assert.Equal(t, account2.Balance, body2.Accounts[1].Balance)
	}
}
