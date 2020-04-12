package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
	"strconv"
	"time"
)

//
// Example shows ability to create a context client for a custom `UUID`
//
func main() {
	util.SetupLogging()
	util.LoadEnv()

	root, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx := root.UUID("my-different-uuid")

	key := fmt.Sprintf("%s", strconv.FormatInt(time.Now().Unix(), 10))
	value := "bar"
	log.Infof("creating key(%s), value(%s)", key, value)
	if err := ctx.Create(key, value, util.GasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}
}
