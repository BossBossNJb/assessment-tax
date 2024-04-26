package tax

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	testCases := []struct {
		name              string
		totalIncome       float64
		wht               float64
		allowances        []Allowance
		personalDeduction float64
		expectedTaxResult float64
		expectedTaxLevels []TaxLevel
	}{
		{
			name:              "Story: EXP01 calculate tax",
			totalIncome:       500000.0,
			wht:               0.0,
			allowances:        []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			expectedTaxResult: 29000.0,
			expectedTaxLevels: []TaxLevel{},
		},
		{
			name:              "Story: EXP02 calculate tax with WHT",
			totalIncome:       500000.0,
			wht:               25000.0,
			allowances:        []Allowance{{AllowanceType: "donation", Amount: 0.0}},
			expectedTaxResult: 4000.0,
			expectedTaxLevels: []TaxLevel{},
		},
		{
			name:              "Story: EXP03 calculate tax with donation reduce",
			totalIncome:       500000.0,
			wht:               0.0,
			allowances:        []Allowance{{AllowanceType: "donation", Amount: 200000.0}},
			expectedTaxResult: 19000.0,
			expectedTaxLevels: []TaxLevel{},
		},
		{
			name:              "Story: EXP04 calculate tax with tax level detail",
			totalIncome:       500000.0,
			wht:               0.0,
			allowances:        []Allowance{{AllowanceType: "donation", Amount: 200000.0}},
			expectedTaxResult: 19000.0,
			expectedTaxLevels: []TaxLevel{
				{"0-150,000", 0.0},
				{"150,001-500,000", 19000.0},
				{"500,001-1,000,000", 0.0},
				{"1,000,001-2,000,000", 0.0},
				{"2,000,001 ขึ้นไป", 0.0},
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the calculateTax function
			response, err := CalculateTax(tc.totalIncome, tc.wht, tc.allowances, tc.personalDeduction)
			if err != nil {
				t.Fatalf("error calculating tax: %v", err)
			}
			actualTaxResult := response.Tax           // Access the Tax of CalculationResponse
			actualTaxLevelResult := response.TaxLevel // Access the TaxLevel  of CalculationResponse

			// Compare the actual result with the expected result
			if actualTaxResult != tc.expectedTaxResult {
				t.Errorf("test case %s: expected result %f; got %f", tc.name, tc.expectedTaxResult, actualTaxResult)
			}

			// Compare the actual tax levels with the expected tax levels
			for i, expectedLevel := range tc.expectedTaxLevels {
				if actualTaxLevelResult[i].Level != expectedLevel.Level || actualTaxLevelResult[i].Tax != expectedLevel.Tax {
					t.Errorf("test case %s: expected tax level %d: %v; got %v", tc.name, i, expectedLevel, actualTaxLevelResult[i])
				}
			}
		})
	}
}
