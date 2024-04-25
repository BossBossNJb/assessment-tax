// tax_handler.go

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

	taxAmount := CalculateTax(request.TotalIncome, request.Allowances)

	response := CalculationResponse{Tax: taxAmount}
	return c.JSON(http.StatusOK, response)
}
