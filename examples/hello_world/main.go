package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
	"strconv"
	"time"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	ctx, err := util.NewClient()
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
	if err := ctx.Create(key, value, util.GasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}

	// read key
	if v, err := ctx.Read(key, false); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("val for key(%s): %s", key, v)
	}
}
