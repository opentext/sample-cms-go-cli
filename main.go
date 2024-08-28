package main

import (
	"ocp/sample/planets/cmd"

	"github.com/joho/godotenv"
)

func main() {
	initConfig()
	cmd.Execute()
}

func initConfig() {
	// Don't throw an error here because we could be loading config directly from 
	// environment variables.
	godotenv.Load()
}
