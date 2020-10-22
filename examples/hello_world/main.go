package main

import (
	"strconv"
	"time"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// read account
	if account, err := ctx.Account(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("account info: %+v", account)
	}

	key := strconv.FormatInt(time.Now().Unix(), 10)
	value := "bar"

	// create key
	if err := ctx.Create(key, value, bluzelle.GetTestGasInfo(), nil); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}

	// read key
	if v, err := ctx.Read(key); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("val for key(%s): %s", key, v)
	}
}
