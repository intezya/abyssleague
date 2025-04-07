package dotenv

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() {
	envFile := flag.String("env-file", ".env", "Path to .env file")
	flag.Parse()
	err := godotenv.Load(*envFile)
	if err != nil {
		log.Printf(
			"Warning: Error loading .env file from %s: %v. Is it specified correctly? use -env-file=... flag",
			*envFile,
			err,
		)
	} else {
		log.Printf("Environment variables loaded from %s", *envFile)
	}
}
