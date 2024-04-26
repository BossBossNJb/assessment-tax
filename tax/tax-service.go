package tax

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
	KreceiptLimitDeduction float64 `json:"kreceiptLimitDeduction"`
}

// TaxDetailsResponse represents the response structure for tax details.
type TaxDetailsResponse struct {
	PersonalDeduction      float64 `json:"personalDeduction"`
	KreceiptLimitDeduction float64 `json:"kreceiptLimitDeduction"`
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
	if request.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, "Amount exceeds the maximum allowed limit")
	}

	// Update the Kreceipt limit deduction value
	KreceiptLimitDeduction = request.Amount

	response := KreceiptLimitDeductionResponse{KreceiptLimitDeduction: KreceiptLimitDeduction}
	return c.JSON(http.StatusOK, response)
}

func TestSetKreceipLimitDeductionHandler(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "Valid",
			requestBody:        `{"amount": 80000}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"kreceiptLimitDeduction":80000}`,
		},
		{
			name:               "ExceedsLimit",
			requestBody:        `{"amount": 150000}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"Amount exceeds the maximum allowed limit"`,
		},
		{
			name:               "InvalidRequest",
			requestBody:        `{"invalid": "data"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"Invalid request"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/set-kreceip-limit-deduction", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := SetKreceipLimitDeductionHandler(c)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(rec.Body.String()))
			assert.NoError(t, err)
		})
	}
}
