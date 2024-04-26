package tax

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	testCases := []struct {
		name           string
		totalIncome    float64
		wht            float64
		allowances     []Allowance
		expectedResult float64
	}{{
		name:        "Story: EXP01 calculate tax",
		totalIncome: 500000.0,
		wht:         0.0,
		allowances: []Allowance{
			{
				AllowanceType: "donation",
				Amount:        0.0,
			},
		},
		expectedResult: 29000.0,
	},
		{
			name:        "Story: EXP02 calculate tax with WHT",
			totalIncome: 500000.0,
			wht:         25000.0,
			allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
			expectedResult: 4000.0,
		},
		{
			name:        "Story: EXP03 calculate tax with donation reduce",
			totalIncome: 500000.0,
			wht:         0.0,
			allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000.0,
				},
			},
			expectedResult: 19000.0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the calculateTax function
			actualResult := CalculateTax(tc.totalIncome, tc.wht, tc.allowances)

			// Compare the actual result with the expected result
			if actualResult != tc.expectedResult {
				t.Errorf("test case %s: expected result %f; got %f", tc.name, tc.expectedResult, actualResult)
			}
		})
	}
}
