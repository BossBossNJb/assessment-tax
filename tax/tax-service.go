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

// AdminDeductionResponse by admin.
type AdminPersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

// KreceiptLimitDeductionResponse response by admin.
type KreceiptLimitDeductionResponse struct {
	// KreceiptLimitDeduction float64 `json:"kreceiptLimitDeduction"`
	KreceiptLimitDeduction float64 `json:"kReceipt"`
}

// TaxDetailsResponse represents the response structure for tax details.
type TaxDetailsResponse struct {
	PersonalDeduction      float64 `json:"personalDeduction"`
	KreceiptLimitDeduction float64 `json:"kReceipt"`
	// KReceipt float64 `json:"kReceipt"`
}

// PersonalDeduction Default .
var PersonalDeduction float64 = 60000.0

// PersonalDeduction Default .
var KreceiptLimitDeduction float64 = 50000.0

// CalculateTaxHandler handles the HTTP request for tax calculation.
func CalculateTaxHandler(c echo.Context) error {
	var request CalculationRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	var donationDeduction float64
	var kreceiptDeduction float64
	for _, allowance := range request.Allowances {
		switch allowance.AllowanceType {
		case "donation":
			donationDeduction = allowance.Amount
		case "k-receipt":
			kreceiptDeduction = allowance.Amount
		}
	}

	// Check for negative values of PersonalDeduction, donation deduction, and k-receipt deduction
	if PersonalDeduction < 0 || donationDeduction < 0 || kreceiptDeduction < 0 {
		return c.JSON(http.StatusBadRequest, "Invalid values for deductions: PersonalDeduction, donation, or k-receipt")
	}

	// Check for WHT is non-negative and does not exceed total income
	if request.WHT < 0 || request.WHT > request.TotalIncome {
		return c.JSON(http.StatusBadRequest, "Invalid value for WHT: must be non-negative and not exceed total income")
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
	if request.Amount < 10000.0 || request.Amount > 100000.0 {
		return c.JSON(http.StatusBadRequest, "Amount exceeds PersonalDeduction the allowed limit")
	}

	// Update the personal deduction value
	PersonalDeduction = request.Amount

	response := AdminPersonalDeductionResponse{PersonalDeduction: PersonalDeduction}
	return c.JSON(http.StatusOK, response)
}

// / TaxDetails handles the HTTP request for tax details.
func TaxDetails(c echo.Context) error {
	response := TaxDetailsResponse{
		PersonalDeduction:      PersonalDeduction,
		KreceiptLimitDeduction: KreceiptLimitDeduction,
	}
	return c.JSON(http.StatusOK, response)
}

// Set KreceipLimitDeductionHandler handles the HTTP request for setting the K-receipt limit deduction by admin.
func SetKreceipLimitDeductionHandler(c echo.Context) error {
	var request AdminDeductionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// Check if the requested amount is within the allowed range
	if request.Amount < 10000.0 || request.Amount > 100000.0 {
		return c.JSON(http.StatusBadRequest, "Amount exceeds KreceipLimitDeduction the allowed limit")
	}

	// Update the Kreceipt limit deduction value
	KreceiptLimitDeduction = request.Amount

	response := KreceiptLimitDeductionResponse{KreceiptLimitDeduction: KreceiptLimitDeduction}

	return c.JSON(http.StatusOK, response)
}
