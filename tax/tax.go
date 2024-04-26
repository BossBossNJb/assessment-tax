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

// TaxLevel represents the tax level structure for tax calculation.
type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

// CalculationResponse represents the response structure for tax calculation.
type CalculationResponse struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}

// calculateTax calculates the tax based on income and allowances.
func CalculateTax(income float64, wht float64, allowances []Allowance) CalculationResponse {
	var tax float64
	var taxFinalPaid float64
	var donationDeduction float64

	// personalAllowance represents the fixed personal allowance.
	personalDeduction := 60000.0

	// Ensure that personal allowance more equal 60000
	if personalDeduction < 600000 {
		personalDeduction = 60000
	}

	// Define tax levels
	taxLevels := []TaxLevel{
		{"0-150,000", 0.0},
		{"150,001-500,000", 0.0},
		{"500,001-1,000,000", 0.0},
		{"1,000,001-2,000,000", 0.0},
		{"2,000,001 ขึ้นไป", 0.0},
	}

	// Calculate donation deduction
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			if allowance.Amount > 100000 { // Ensure that donation allowance limit is 100000
				donationDeduction = 100000
			} else if allowance.Amount < 0 { // Ensure that donation allowance is not negative
				donationDeduction = 0
			} else {
				donationDeduction = allowance.Amount
			}
			break
		}
	}

	// Calculate taxable income after deductions
	incomeAfterDeductions := income - personalDeduction - donationDeduction

	// Ensure that income after deductions is not negative
	if incomeAfterDeductions < 0 {
		incomeAfterDeductions = 0
	}
	taxableIncome := incomeAfterDeductions

	// Calculate tax for each level
	for i, _ := range taxLevels {
		tax = 0
		switch i {
		case 0: // "0-150,000"
			if taxableIncome < 150000 {
				tax = 0
			}
		case 1: // "150,001-500,000"
			if taxableIncome > 150000 {
				if taxableIncome > 500000 {
					tax = 35000
				} else {
					tax = (taxableIncome - 150000) * 0.10
				}
			}
		case 2: // "500,001-1,000,000"
			if taxableIncome > 500000 {
				if taxableIncome > 1000000 {
					tax = 95000
				} else {
					tax = (taxableIncome - 500000) * 0.15
				}
			}
		case 3: // "1,000,001-2,000,000"
			if taxableIncome > 1000000 {
				if taxableIncome > 2000000 {
					tax = 335000
				} else {
					tax = (taxableIncome - 1000000) * 0.20
				}
			}
		case 4: // "2,000,001 more"
			if taxableIncome > 2000000 {
				tax = (taxableIncome - 2000000) * 0.35
			}
		}
		taxLevels[i].Tax = tax
	}

	// Calculate tax total from sum tax levels
	taxTotal := 0.0
	for _, level := range taxLevels {
		taxTotal += level.Tax
	}

	// Ensure if withholding tax exceeds the limit 100000
	if wht > 100000 {
		wht = 100000
	}

	// Calculate tax final paid on taxable income after deductions including withholding tax
	taxFinalPaid = taxTotal - wht

	// Ensure tax is not negative
	if taxFinalPaid < 0 {
		taxFinalPaid = 0
	}

	// Return the tax value from the CalculationResponse instance
	return CalculationResponse{
		Tax:      taxFinalPaid,
		TaxLevel: taxLevels,
	}
}
