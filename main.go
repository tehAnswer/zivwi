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
	e.POST("/authorize", authorize)
	e.POST("/cmd/pay", payCmd)
	e.GET("/q/balance", balanceQuery)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}

// Handlers
func authorize(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented.")
}

func payCmd(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented.")
}

func balanceQuery(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "Not implemented")
}
