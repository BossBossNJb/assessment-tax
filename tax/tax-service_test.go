package tax

import (
	"bytes"
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
		"allowances": [
			{
				"allowanceType": "donation",
				"amount": 0.0
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
	expectedResponseBody := `{"tax":29000}`
	actualResponseBody := strings.TrimSpace(rec.Body.String()) // Trim the newline character
	assert.Equal(t, expectedResponseBody, actualResponseBody)

}
