package main

import (
	"fmt"
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

	no_of_creates := 4
	creates := make(chan Result, no_of_creates)

	go withTransaction(ctx, "with txn(1)", creates)
	go withTransaction(ctx, "with txn(2)", creates)
	go withTransaction(ctx, "with txn(3)", creates)
	go withoutTransaction(ctx, "without txn", creates)

	i := 0
	for ret := range creates {
		if ret.Error != nil {
			log.Fatalf("error(%s): %s", ret.Type, err)
		} else {
			log.Infof("done(%s)", ret.Type)
		}
		i++
		if i == no_of_creates {
			close(creates)
		}
	}
}

func withTransaction(ctx *bluzelle.Client, typ string, c chan Result) {
	txn := ctx.Transaction()
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%s-%s", typ, makeRandomString(5))
		log.Infof("with txn: init create key(%s)", key)
		if err := txn.Create(key, "value", bluzelle.GetTestGasInfo(), nil); err != nil {
			c <- Result{
				err, typ,
			}
			return
		}
	}
	err := txn.Execute()
	c <- Result{
		err, typ,
	}
}

func withoutTransaction(ctx *bluzelle.Client, typ string, c chan Result) {
	key := fmt.Sprintf("%s-%s", typ, makeRandomString(5))
	log.Infof("without txn: init create key(%s)", key)
	err := ctx.Create(key, "value", bluzelle.GetTestGasInfo(), nil)
	c <- Result{
		err, typ,
	}
}

func makeRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
