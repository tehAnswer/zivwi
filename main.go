package main

import (
	"flag"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	address = flag.String("port", "127.0.0.1:3000", "TCP address to listen to.")
)

func main() {
	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/authorize", Authorize)
	e.POST("/cmd/pay", PayCmd)
	e.GET("/q/balance", BalanceQuery)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}

// Handlers
func Authorize(c echo.Context) error {
	credentials := new(struct {
		Email    string `json:email`
		Password string `json:password`
	})

	if err := c.Bind(credentials); err != nil {
		return c.JSON(422, struct {
			Message string `json:message`
		}{"Incorrect login params, please send email and password"})
	}

	token, loginErr := NewAuthorizeService().Login(credentials.Email, credentials.Password)
	if loginErr != nil {
		return c.String(401, loginErr.Error())
	}

	return c.JSON(200, struct {
		Token string `json:"token"`
	}{token})
}

func PayCmd(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented.")
}

func BalanceQuery(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}
