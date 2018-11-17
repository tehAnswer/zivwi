package main

import (
	"flag"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
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

	// Routes
	e.POST("/authorize", Authorize)

	commands := e.Group("/cmd")
	queries := e.Group("/q")

	commands.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
	queries.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))

	commands.POST("/transfer", TransferCmd)
	queries.GET("/balance", BalanceQuery)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}

// Handlers
func Authorize(c echo.Context) error {
	credentials := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.Bind(credentials); err != nil {
		return c.JSON(422, struct {
			Message string `json:"message"`
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

func TransferCmd(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountIds := claims["account_ids"].([]interface{})

	params := new(struct {
		FromAccountId string `json:"from_account_id"`
		ToAccountId   string `json:"to_account_id"`
		Message       string `json:"message"`
		Amount        uint64 `json:"amount"`
	})

	if paramsErr := c.Bind(params); paramsErr != nil {
		return c.JSON(422, struct {
			Message string `json:"message"`
		}{"Unprocessable transfer."})
	}

	if !contains(accountIds, params.FromAccountId) {
		return c.JSON(401, struct {
			Message string `json:"message"`
		}{"Unauthorized access."})
	}

	transfer, err := appCtx.TransferService.Perform(
		params.FromAccountId,
		params.ToAccountId,
		params.Amount,
		params.Message)

	if err != nil {
		return c.JSON(422, struct {
			Message string `json:message`
		}{err.Error()})
	}

	return c.JSON(201, transfer)
}

func BalanceQuery(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}

func contains(s []interface{}, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
