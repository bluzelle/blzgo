package main

import (
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo/examples/util"
	"github.com/vbstreetz/blzgo"
	// "os"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	// args := os.Args[1:]
	// if len(args) == 0 {
	// 	log.Fatalf("at least one key=value pair is required")
	// }

	ctx, err := util.NewClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	keyValues := []*bluzelle.KeyValue{}
	keyValues = append(keyValues, &bluzelle.KeyValue{"foo", "bar"})
	keyValues = append(keyValues, &bluzelle.KeyValue{"1", "2"})

	log.Infof("updating keys(%s)...", keyValues)

	if err := ctx.MultiUpdate(keyValues); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("done")
	}
}
