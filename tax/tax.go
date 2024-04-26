package tax

import (
	"errors"
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
func CalculateTax(income float64, wht float64, allowances []Allowance, personalDeduction float64) (CalculationResponse, error) {
	var tax float64
	var taxFinalPaid float64
	var donationDeduction float64
	var kreceiptDeduction float64
	// Define tax levels
	taxLevels := []TaxLevel{
		{"0-150,000", 0.0},
		{"150,001-500,000", 0.0},
		{"500,001-1,000,000", 0.0},
		{"1,000,001-2,000,000", 0.0},
		{"2,000,001 ขึ้นไป", 0.0},
	}

	// personalAllowance represents the fixed personal allowance.
	if personalDeduction < 0 { // Ensure that personal deductio is not negative
		personalDeduction = 0
	} else if personalDeduction < 60000 { // Ensure that personal deduction is at least 60000
		personalDeduction = 60000
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
		}

		if allowance.AllowanceType == "k-receipt" {
			if allowance.Amount > 50000 { // Ensure that kreceipt allowance limit is 100000
				kreceiptDeduction = 50000
			} else if allowance.Amount < 0 { // Ensure that kreceipt allowance is not negative
				kreceiptDeduction = 0
			} else {
				kreceiptDeduction = allowance.Amount
			}
		}
	}

	// Calculate taxable income after deductions
	incomeAfterDeductions := income - personalDeduction - donationDeduction - kreceiptDeduction

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

	// Ensure that the withholding tax provided for the calculation does not exceed the total income.
	if wht > income {
		return CalculationResponse{}, errors.New("withholding tax cannot be greater than the total income")
	}

	// withholding represents the fixed personal allowance.
	if wht < 0 { // Ensure that withholding  is not negative
		wht = 0
	} else if wht > 100000 { // Ensure if withholding tax exceeds the limit 100000
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
	}, nil
}
