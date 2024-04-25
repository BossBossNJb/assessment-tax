package main

import (
	"net/http"

	"github.com/BossBossNJb/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Root endpoint handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	// Tax calculation endpoint handler
	e.POST("/tax/calculations", tax.CalculateTaxHandler)

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
