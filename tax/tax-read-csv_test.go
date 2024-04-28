package tax

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestReadCSVHandler(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create a sample CSV content
	csvContent := `totalIncome,wht,donation
	500000,0,0
	600000,40000,20000
	750000,50000,15000`

	// Create a temporary file to store the CSV content
	file, err := ioutil.TempFile("", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // Clean up the temporary file

	// Write the CSV content to the temporary file
	if _, err := file.WriteString(csvContent); err != nil {
		t.Fatal(err)
	}

	// Close the file
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	// Open the temporary file for reading
	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a multipart form writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form file part
	part, err := writer.CreateFormFile("taxFile", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	// Copy the file content to the form file part
	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a request with the multipart form data
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a response recorder
	rec := httptest.NewRecorder()

	// Create a context with the request and response recorder
	c := e.NewContext(req, rec)

	// Call the ReadCSVHandler function
	records, err := ReadCSVHandler(c)

	// Check if there's no error
	assert.NoError(t, err)

	// Check the number of records
	assert.Equal(t, 4, len(records)) // Excluding header row
}

func TestCalculateTaxFromCSV(t *testing.T) {
	// Sample records
	records := [][]string{
		{"totalIncome", "wht", "donation"},
		{"500000", "0", "0"},
		{"600000", "40000", "20000"},
		{"750000", "50000", "15000"},
	}

	// Call the CalculateTaxFromCSV function
	taxCalculations, err := CalculateTaxFromCSV(records)

	// Check if there's no error
	assert.NoError(t, err)

	// Check the number of tax calculations
	assert.Equal(t, 3, len(taxCalculations))
}

func TestCalculateTaxFromCSVHandler(t *testing.T) {
	// Create a new echo instance
	e := echo.New()

	// Create a sample CSV content
	csvContent := `totalIncome,wht,donation
	500000,0,0
	600000,40000,20000
	750000,50000,15000`

	// Create a temporary file to store the CSV content
	file, err := ioutil.TempFile("", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // Clean up the temporary file

	// Write the CSV content to the temporary file
	if _, err := file.WriteString(csvContent); err != nil {
		t.Fatal(err)
	}

	// Close the file
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	// Open the temporary file for reading
	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a multipart form writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form file part
	part, err := writer.CreateFormFile("taxFile", filepath.Base(file.Name()))
	if err != nil {
		t.Fatal(err)
	}

	// Copy the file content to the form file part
	if _, err := io.Copy(part, file); err != nil {
		t.Fatal(err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	// Create a request with the multipart form data
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a response recorder
	rec := httptest.NewRecorder()

	// Create a context with the request and response recorder
	c := e.NewContext(req, rec)

	// Call the CalculateTaxFromCSVHandler function
	err = CalculateTaxFromCSVHandler(c)

	// Check if there's no error
	assert.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the response body
	expectedResponseBody := `{"taxes":[{"totalIncome":500000,"tax":29000},{"totalIncome":600000,"tax":0},{"totalIncome":750000,"tax":11250}]}`
	assert.Equal(t, expectedResponseBody, strings.TrimSpace(rec.Body.String()))
}
