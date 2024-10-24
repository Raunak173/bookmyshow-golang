package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	// Access environment variables
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	serviceId := os.Getenv("TWILIO_SERVICE_SID")

	// Print to check if variables are loaded correctly
	fmt.Println("Account SID:", accountSID)
	fmt.Println("Auth Token:", authToken)
	fmt.Println("Service SID:", serviceId)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
