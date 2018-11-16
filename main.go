package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	address = flag.String("port", "127.0.0.1:3000", "TCP address to listen to.")
	appCtx  = NewAppCtx()
)

func main() {
	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))

	// Routes
	e.POST("/authorize", Authorize)

	commands := e.Group("/cmd")
	queries := e.Group("/q")

	commands.POST("/pay", PayCmd)
	queries.GET("/balance", BalanceQuery)

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

	token, loginErr := appCtx.AuthorizeService.Login(credentials.Email, credentials.Password)
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
