package tax

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AdminDeductionRequest represents the request structure for setting personal deduction by admin.
type AdminDeductionRequest struct {
	Amount float64 `json:"amount"`
}

// AdminDeductionResponse represents the response structure for setting personal deduction by admin.
type AdminDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

// PersonalDeduction Default .
var PersonalDeduction float64 = 60000.0

// CalculateTaxHandler handles the HTTP request for tax calculation.
func CalculateTaxHandler(c echo.Context) error {
	var request CalculationRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// Calculate tax amount and tax levels
	response, err := CalculateTax(request.TotalIncome, request.WHT, request.Allowances, PersonalDeduction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error calculating tax: %v", err))
	}

	// Return the response
	return c.JSON(http.StatusOK, response)
}

// SetPersonalDeductionHandler handles the HTTP request for setting personal deduction by admin.
func SetPersonalDeductionHandler(c echo.Context) error {
	var request AdminDeductionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// Check if the requested amount is within the allowed range
	if request.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, "Amount exceeds the maximum allowed limit")
	}

	// Update the personal deduction value
	PersonalDeduction = request.Amount

	response := AdminDeductionResponse{PersonalDeduction: PersonalDeduction}
	return c.JSON(http.StatusOK, response)
}

// TaxDetails with PersonalDeduction
func TaxDetails(c echo.Context) error {
	response := AdminDeductionResponse{PersonalDeduction: PersonalDeduction}
	return c.JSON(http.StatusOK, response)
}
