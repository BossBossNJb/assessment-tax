package tax

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// TaxData represents tax-related data from the CSV file
type TaxData struct {
	TotalIncome float64 `csv:"totalIncome"`
	WHT         float64 `csv:"wht"`
	Donation    float64 `csv:"donation"`
}

// TaxCalculation represents the calculated tax for a set of tax data
type TaxCalculation struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

// ReadCSVHandler parses the CSV file from the request body
func ReadCSVHandler(c echo.Context) ([][]string, error) {
	// Parse the CSV file from the request body
	file, err := c.FormFile("taxFile")
	if err != nil {
		return nil, err
	}

	// Open the uploaded CSV file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Parse the CSV file
	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// CalculateTaxFromCSV calculates tax from CSV records
func CalculateTaxFromCSV(records [][]string) ([]TaxCalculation, error) {
	var taxCalculations []TaxCalculation
	for _, record := range records {
		// Cheak csv format
		if strings.Contains(record[0], "totalIncome") || strings.Contains(record[0], "wht") || strings.Contains(record[0], "donation") {
			if record[0] == "totalIncome" || record[1] == "wht" || record[2] == "donation" {
				continue
			} else {
				return nil, fmt.Errorf("invalid CSV format")
			}
		}

		// Validate and parse the record
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid CSV format")
		}

		totalIncomeStr := strings.TrimSpace(record[0])
		totalIncomeStr = strings.ReplaceAll(totalIncomeStr, ",", ".")
		totalIncome, err := strconv.ParseFloat(totalIncomeStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid totalIncome")
		}

		wht, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid WHT")
		}

		donation, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid donation")
		}

		// Calculate tax using the existing CalculateTax function
		taxResponse, err := CalculateTax(totalIncome, wht, []Allowance{{AllowanceType: "donation", Amount: donation}}, PersonalDeduction, KreceiptLimitDeduction)
		if err != nil {
			return nil, err
		}

		// Append tax calculation to the response
		taxCalculations = append(taxCalculations, TaxCalculation{
			TotalIncome: totalIncome,
			Tax:         taxResponse.Tax,
		})
	}

	return taxCalculations, nil
}

// CalculateTaxFromCSV handles CSV parsing and tax calculation
func CalculateTaxFromCSVHandler(c echo.Context) error {
	// Parse CSV
	records, err := ReadCSVHandler(c)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing CSV file: %v", err))
	}

	// Calculate tax
	taxCalculations, err := CalculateTaxFromCSV(records)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error calculating tax: %v", err))
	}

	// Return the calculated taxes as JSON
	return c.JSON(http.StatusOK, map[string]interface{}{"taxes": taxCalculations})
}
