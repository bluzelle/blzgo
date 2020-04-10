package main

import (
	"github.com/apex/log"
	clih "github.com/apex/log/handlers/cli"
	jsonh "github.com/apex/log/handlers/json"
	"github.com/joho/godotenv"
	"os"
)

func setupLogging() {
	if os.Getenv("ENV") == "" { // dev
		log.SetHandler(clih.Default)
	} else {
		log.SetHandler(jsonh.Default)
	}
	log.SetLevel(log.DebugLevel)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		// log.Errorf("%s", err)
		if err := godotenv.Load("../.env"); err != nil { // when running tests
			// log.Errorf("%s", err)
		}
	}
}
