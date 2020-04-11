package main

import (
	"github.com/apex/log"
	clih "github.com/apex/log/handlers/cli"
	"github.com/joho/godotenv"
)

func setupLogging() {
	log.SetHandler(clih.Default)
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
