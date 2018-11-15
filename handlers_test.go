package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	main "github.com/tehAnswer/zivwi"
)

var (
	correctCredentials   = `{"email": "benito@rome.it", "password":"cia0p0rc0di0"}`
	incorrectCredentials = `{"email": "benito@rome.it", "password":"ciao"}`
)

func TestLoginThroughAPI(t *testing.T) {
	gateway := main.NewUserGateway()
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
		strings.NewReader(correctCredentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, main.Authorize(c)) {
		assert.Equal(t, 200, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}
}

func TestIncorrectLoginThroughAPI(t *testing.T) {
	gateway := main.NewUserGateway()
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
