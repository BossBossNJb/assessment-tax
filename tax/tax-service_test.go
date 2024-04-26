package tax

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTaxHandler(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create the request payload
	requestBody := `{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 200000.0
		  }
		]
	}`

	// Create an HTTP request with the payload
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Create a context with the request and response recorder
	c := e.NewContext(req, rec)

	// Call the CalculateTaxHandler function
	err := CalculateTaxHandler(c)
	assert.NoError(t, err)

	// Check if the status code is OK
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the response body
	expectedResponseBody := CalculationResponse{
		Tax: 19000,
		TaxLevel: []TaxLevel{
			{"0-150,000", 0},
			{"150,001-500,000", 19000},
			{"500,001-1,000,000", 0},
			{"1,000,001-2,000,000", 0},
			{"2,000,001 ขึ้นไป", 0},
		},
	}
	expectedJSON, err := json.Marshal(expectedResponseBody)
	if err != nil {
		t.Fatalf("failed to marshal expected response body: %v", err)
	}

	// Check the response body
	actualResponseBody := strings.TrimSpace(rec.Body.String()) // Trim the newline character
	assert.Equal(t, string(expectedJSON), actualResponseBody)

}
