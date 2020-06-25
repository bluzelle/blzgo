package main

import (
	"strconv"
	"time"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
	// "os"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	// args := os.Args[1:]
	// if len(args) == 0 {
	// 	log.Fatalf("at least one key=value pair is required")
	// }

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key1 := strconv.FormatInt(100+time.Now().Unix(), 10)
	key2 := strconv.FormatInt(200+time.Now().Unix(), 10)

	if err := ctx.Create(key1, "value", bluzelle.TestGasInfo(), nil); err != nil {
		log.Fatalf("%s", err)
	}
	if err := ctx.Create(key2, "value", bluzelle.TestGasInfo(), nil); err != nil {
		log.Fatalf("%s", err)
	}

	keyValues := []*bluzelle.KeyValue{}
	keyValues = append(keyValues, &bluzelle.KeyValue{key1, "bar"})
	keyValues = append(keyValues, &bluzelle.KeyValue{key2, "baz"})

	log.Infof("updating keys(%s)...", keyValues)

	if err := ctx.MultiUpdate(keyValues, bluzelle.TestGasInfo()); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("done")
	}
}
