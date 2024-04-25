// tax_test.go

package tax

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	// Define the test case using the provided JSON payload
	totalIncome := 500000.0
	allowances := []Allowance{
		{
			AllowanceType: "donation",
			Amount:        0.0,
		},
	}
	expectedTax := 29000.0

	// Call the calculateTax function
	actualTax := CalculateTax(totalIncome, allowances)

	// Compare the actual tax with the expected tax
	if actualTax != expectedTax {
		t.Errorf("expected tax %f; got %f", expectedTax, actualTax)
	}
}
