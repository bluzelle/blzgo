package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo"
	"strconv"
	"time"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// read account
	if account, err := ctx.ReadAccount(); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("account info: %+v", account)
	}

	key := strconv.FormatInt(time.Now().Unix(), 10)
	value := "bar"

	// create key
	if err := ctx.Create(key, value, nil, nil); err != nil {
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
