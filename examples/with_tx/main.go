package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

type Result struct {
	Error error
	Type  string
}

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	txn := ctx.Transaction()

	if err := txn.Create(key(), "value", bluzelle.GetTestGasInfo(), nil); err != nil {
		log.Fatalf("%s", err)
	}

	if err := txn.TxRead(key(), bluzelle.GetTestGasInfo()); err != nil {
		log.Fatalf("%s", err)
	}

	if err := txn.Execute(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("txread result %s", txt.GetTxReadResult(1))
	}
}

func key() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < 10; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
