package tax

import (
	_ "net/http"

	_ "github.com/labstack/echo/v4"
)

// Allowance represents a type of allowance with its amount.
type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

// CalculationRequest represents the request structure for tax calculation.
type CalculationRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

// CalculationResponse represents the response structure for tax calculation.
type CalculationResponse struct {
	Tax float64 `json:"tax"`
}

// calculateTax calculates the tax based on income and allowances.
func CalculateTax(income float64, wht float64, allowances []Allowance) float64 {
	personalAllowance := 60000.0
	var tax float64

	// Calculate taxable income after deductions
	incomeAfterDeductions := income - personalAllowance
	taxableIncome := incomeAfterDeductions

	// Calculate tax based on taxable income
	if taxableIncome <= 150000 {
		tax = 0
	} else if taxableIncome <= 500000 {
		tax = (taxableIncome - 150000) * 0.1
	} else if taxableIncome <= 1000000 {
		tax = 35000 + (taxableIncome-500000)*0.15
	} else if taxableIncome <= 2000000 {
		tax = 95000 + (taxableIncome-1000000)*0.2
	} else {
		tax = 335000 + (taxableIncome-2000000)*0.35
	}

	taxFinalPaid := tax - wht

	// Ensure tax is not negative
	if taxFinalPaid < 0 {
		taxFinalPaid = 0
	}

	return taxFinalPaid
}
