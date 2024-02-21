package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string
	for _, fieldError := range err.(validator.ValidationErrors) {
		errors = append(errors, fieldError.Error())
	}
	return errors
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
