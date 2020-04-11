package main

import (
	"fmt"
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

	for value := 0; value < 10; value++ {
		t := strconv.FormatInt(time.Now().Unix(), 10)
		key := fmt.Sprintf("%s-%d", t, value)
		log.Infof("creating key(%s), value(%d)", key, value)
		if err := ctx.Create(key, fmt.Sprintf("%d", value), util.GasInfo()); err != nil {
			log.Fatalf("%s", err)
		} else {
			log.Infof("created key")
		}
	}
}
