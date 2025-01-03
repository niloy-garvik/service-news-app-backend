package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() error {
	if err := godotenv.Load(); err != nil {
		return errors.New("no env file found")
	} else {
		return nil
	}
}

func GetEnvironmentVariable(variableName string) string {

	variableValue := os.Getenv(variableName)
	if variableValue == "" {
		// print("\nno such env variable found named: ", variableName)
		return ""
	} else {
		return variableValue
	}
}
