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

	// Define the route for setting personal deduction by admin
	e.POST("/admin/deductions/personal", tax.SetPersonalDeductionHandler) // Update the route handler

	// Tax calculation endpoint handler
	e.POST("/tax/calculations", tax.CalculateTaxHandler)
	e.GET("/tax/calculations", tax.TaxDetails)

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
