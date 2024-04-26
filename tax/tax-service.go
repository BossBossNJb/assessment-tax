package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CalculateTaxHandler handles the HTTP request for tax calculation.
func CalculateTaxHandler(c echo.Context) error {
	var request CalculationRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// Calculate tax amount and tax levels
	response := CalculateTax(request.TotalIncome, request.WHT, request.Allowances)

	// Return the response
	return c.JSON(http.StatusOK, response)
}
