package main

import (
	"fmt"
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
	if err := ctx.Create(key, value, bluzelle.TestGasInfo(), nil); err != nil {
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

func withTransaction(ctx *bluzelle.Client) {
	txn := ctx.Transaction()
	for i := 0; i < 10; i++ {
		if err := txn.Create(fmt.Sprintf("%d", i), "value", bluzelle.TestGasInfo(), nil); err != nil {
			log.Fatalf("%s", err)
		} else {
			log.Infof("created key")
		}
	}
	if err := txn.Execute(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}
}

func withoutTransaction(ctx *bluzelle.Client) {
	for i := 0; i < 10; i++ {
		if err := ctx.Create(fmt.Sprintf("%d", i), "value", bluzelle.TestGasInfo(), nil); err != nil {
			log.Fatalf("%s", err)
		} else {
			log.Infof("created key")
		}
	}
}
