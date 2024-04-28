package tax

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCalculateTaxHandler(t *testing.T) {
// 	// Define test cases
// 	testCases := []struct {
// 		name               string
// 		requestBody        string
// 		expectedStatusCode int
// 		expectedTaxResult  float64
// 		expectedTaxLevels  []TaxLevel
// 	}{
// 		{
// 			name: "CalculateTax with donation",
// 			requestBody: `{
// 				"totalIncome": 500000.0,
// 				"wht": 0.0,
// 				"allowances": [
// 				  {
// 					"allowanceType": "donation",
// 					"amount": 200000.0
// 				  }
// 				]
// 			}`,
// 			expectedStatusCode: http.StatusOK,
// 			expectedTaxResult:  19000,
// 			expectedTaxLevels: []TaxLevel{
// 				{"0-150,000", 0},
// 				{"150,001-500,000", 19000},
// 				{"500,001-1,000,000", 0},
// 				{"1,000,001-2,000,000", 0},
// 				{"2,000,001 ขึ้นไป", 0},
// 			},
// 		},
// 		{
// 			name: "Check WHT exceeding total income",
// 			requestBody: `{
// 				"totalIncome": 500000.0,
// 				"wht": 550000.0,
// 				"allowances": [
// 				  {
// 					"allowanceType": "donation",
// 					"amount": 200000.0
// 				  }
// 				]
// 			}`,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedTaxResult:  0,
// 			expectedTaxLevels:  nil, // No tax levels expected in case of error
// 		},
// 		{
// 			name: "Check WHT is negative",
// 			requestBody: `{
// 				"totalIncome": 1000.0,
// 				"wht": -550000.0,
// 				"allowances": [
// 				  {
// 					"allowanceType": "donation",
// 					"amount": 200000.0
// 				  }
// 				]
// 			}`,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedTaxResult:  0,
// 			expectedTaxLevels:  nil, // No tax levels expected in case of error
// 		},
// 		{
// 			name: "Check donation is negative",
// 			requestBody: `{
// 				"totalIncome": 500000.0,
// 				"wht": 50000.0,
// 				"allowances": [
// 					{
// 						"allowanceType": "donation",
// 						"amount": -200000.0
// 					}
// 				]
// 			}`,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedTaxResult:  0,
// 			expectedTaxLevels:  nil, // No tax levels expected in case of error
// 		},
// 		{
// 			name: "Check donation is negative",
// 			requestBody: `{
// 				"totalIncome": 500000.0,
// 				"wht": 50000.0,
// 				"allowances": [
// 					{
// 						"allowanceType": "k-receipt",
// 						"amount": -200000.0
// 					}
// 				]
// 			}`,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedTaxResult:  0,
// 			expectedTaxLevels:  nil, // No tax levels expected in case of error
// 		},
// 	}

// 	e := echo.New()

// 	// Run test cases
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBufferString(tc.requestBody))
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 			rec := httptest.NewRecorder()

// 			c := e.NewContext(req, rec)

// 			// Call CalculateTaxHandler function
// 			err := CalculateTaxHandler(c)
// 			assert.NoError(t, err)

// 			// Check status code
// 			assert.Equal(t, tc.expectedStatusCode, rec.Code)

// 			// Check response body if status code is OK
// 			if tc.expectedStatusCode == http.StatusOK {
// 				var response CalculationResponse
// 				err = json.Unmarshal(rec.Body.Bytes(), &response)
// 				assert.NoError(t, err)

// 				// Check tax result
// 				assert.Equal(t, tc.expectedTaxResult, response.Tax)

// 				// Check tax levels
// 				assert.Equal(t, len(tc.expectedTaxLevels), len(response.TaxLevel))
// 				for i, expectedLevel := range tc.expectedTaxLevels {
// 					assert.Equal(t, expectedLevel.Level, response.TaxLevel[i].Level)
// 					assert.Equal(t, expectedLevel.Tax, response.TaxLevel[i].Tax)
// 				}
// 			}
// 		})
// 	}
// }

// func TestSetPersonalDeductionHandler(t *testing.T) {
// 	testCases := []struct {
// 		name               string
// 		request            AdminDeductionRequest
// 		expectedStatusCode int
// 		expectedResponse   string // Expected response body
// 	}{
// 		{
// 			name: "PersonalDeductionValidAmount",
// 			request: AdminDeductionRequest{
// 				Amount: 70000,
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   `{"personalDeduction":70000}`,
// 		},
// 		{
// 			name: "PersonalDeductionExceedingLimitUpper",
// 			request: AdminDeductionRequest{
// 				Amount: 105000,
// 			},
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   `"Amount exceeds PersonalDeduction the allowed limit"`,
// 		},
// 		{
// 			name: "PersonalDeductionExceedingLimitLower",
// 			request: AdminDeductionRequest{
// 				Amount: 5000,
// 			},
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   `"Amount exceeds PersonalDeduction the allowed limit"`,
// 		},
// 	}
// 	// Run test cases
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			e := echo.New()

// 			requestBody, err := json.Marshal(tc.request)
// 			if err != nil {
// 				t.Fatalf("failed to marshal request body: %v", err)
// 			}

// 			req := httptest.NewRequest(http.MethodPost, "/set-personal-deduction", bytes.NewBuffer(requestBody))
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 			rec := httptest.NewRecorder()

// 			c := e.NewContext(req, rec)

// 			// Call SetPersonalDeductionHandler function
// 			err = SetPersonalDeductionHandler(c)

// 			// Check status code matches
// 			assert.Equal(t, tc.expectedStatusCode, rec.Code)

// 			// Check response body
// 			actualResponseBody := strings.TrimSpace(rec.Body.String())
// 			assert.Equal(t, tc.expectedResponse, actualResponseBody)

// 			// Check error
// 			assert.NoError(t, err)
// 		})
// 	}

// }

// func TestSetKreceipLimitDeductionHandler(t *testing.T) {
// 	testCases := []struct {
// 		name               string
// 		request            AdminDeductionRequest
// 		expectedStatusCode int
// 		expectedResponse   string // Expected response body
// 	}{
// 		{
// 			name: "SetKreceipLimitDeductionValidAmount",
// 			request: AdminDeductionRequest{
// 				Amount: 70000,
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   `{"kReceipt":70000}`,
// 		},
// 		{
// 			name: "SetKreceipLimitDeductionExceedingLimitUpper",
// 			request: AdminDeductionRequest{
// 				Amount: 105000,
// 			},
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   `"Amount exceeds KreceipLimitDeduction the allowed limit"`,
// 		},
// 		{
// 			name: "SetKreceipLimitDeductionExceedingLimitLower",
// 			request: AdminDeductionRequest{
// 				Amount: -1000.0,
// 			},
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   `"Amount exceeds KreceipLimitDeduction the allowed limit"`,
// 		},
// 	}
// 	// Run test cases
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			e := echo.New()

// 			requestBody, err := json.Marshal(tc.request)
// 			if err != nil {
// 				t.Fatalf("failed to marshal request body: %v", err)
// 			}

// 			req := httptest.NewRequest(http.MethodPost, "/set-personal-deduction", bytes.NewBuffer(requestBody))
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 			rec := httptest.NewRecorder()

// 			c := e.NewContext(req, rec)

// 			// Call SetPersonalDeductionHandler function
// 			err = SetKreceipLimitDeductionHandler(c)

// 			// Check status code matches
// 			assert.Equal(t, tc.expectedStatusCode, rec.Code)

// 			// Check response body
// 			actualResponseBody := strings.TrimSpace(rec.Body.String())
// 			assert.Equal(t, tc.expectedResponse, actualResponseBody)

// 			// Check error
// 			assert.NoError(t, err)
// 		})
// 	}

// }

// // TestTaxDetails tests the TaxDetails function.
// func TestTaxDetails(t *testing.T) {
// 	testCases := []struct {
// 		name               string
// 		expectedStatusCode int
// 		expectedResponse   string
// 	}{
// 		{
// 			name:               "Valid",
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   `{"personalDeduction":70000,"kReceipt":70000}`,
// 		},
// 	}
// 	// Run test cases
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			e := echo.New()

// 			req := httptest.NewRequest(http.MethodGet, "/tax/calculations/details", nil)
// 			rec := httptest.NewRecorder()

// 			c := e.NewContext(req, rec)

// 			// Call TaxDetails function
// 			err := TaxDetails(c)

// 			// Check status code matches
// 			assert.Equal(t, tc.expectedStatusCode, rec.Code)

// 			// Check response body
// 			actualResponseBody := strings.TrimSpace(rec.Body.String())
// 			assert.Equal(t, tc.expectedResponse, actualResponseBody)

// 			// Check error
// 			assert.NoError(t, err)
// 		})
// 	}
// }
