package util

import (
	"github.com/apex/log"
	clih "github.com/apex/log/handlers/cli"
	"github.com/joho/godotenv"
	"github.com/vbstreetz/blzgo"
	"os"
	"strconv"
)

func NewClient() (*bluzelle.Client, error) {
	debug := false
	if d, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
		debug = d
	}

	// create client
	options := &bluzelle.Options{
		Address:  os.Getenv("ADDRESS"),
		Mnemonic: os.Getenv("MNEMONIC"),
		UUID:     os.Getenv("UUID"),
		Endpoint: os.Getenv("ENDPOINT"),
		ChainId:  os.Getenv("CHAIN_ID"),
		Debug:    debug,
	}
	ctx, err := bluzelle.NewClient(options)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func GasInfo() *bluzelle.GasInfo {
	return &bluzelle.GasInfo{
		MaxFee: 4000001,
	}
}

func SetupLogging() {
	log.SetHandler(clih.Default)
	log.SetLevel(log.DebugLevel)
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		// log.Errorf("%s", err)
		if err := godotenv.Load("../.env"); err != nil { // when running tests
			// log.Errorf("%s", err)
		}
	}
}
