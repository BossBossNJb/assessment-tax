package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BossBossNJb/assessment-tax/tax"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// envFilePath := ".env"
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the values of environment variables
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	// databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	// Root endpoint handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	// Define a custom middleware for admin API group
	adminAuthMiddleware := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Check if the provided username and password match the admin credentials
		if username == adminUsername && password == adminPassword {
			return true, nil
		}
		return false, nil
	})

	// Group the admin API routes and apply basic authentication middleware
	adminGroup := e.Group("/admin")
	adminGroup.Use(adminAuthMiddleware)

	// Define the route for setting personal deduction by admin
	adminGroup.POST("/deductions/personal", tax.SetPersonalDeductionHandler)

	// Define the route for setting k-receipt limit deduction by admin
	adminGroup.POST("/deductions/k-receipt", tax.SetKreceipLimitDeductionHandler)

	// Group tax-related endpoints
	taxGroup := e.Group("/tax")

	// Tax calculation endpoint handler
	taxGroup.POST("/calculations", tax.CalculateTaxHandler)
	taxGroup.GET("/calculations/deteils", tax.TaxDetails)

	// Start the server
	// e.Logger.Fatal(e.Start(":" + port))
	fmt.Println("port:", port)
	fmt.Println("adminUsername:", adminUsername)
	e.Logger.Fatal(e.Start(":" + port))

	// Graceful Shutdown
	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt)
	// <-shutdown
	// // Print "shutting down the server"
	// fmt.Println("Shutting down the server...")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// if err := e.Shutdown(ctx); err != nil {
	// 	e.Logger.Fatal(err)
	// }
}
