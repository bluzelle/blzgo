package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

//
// Example tests correct sequence increment
// when working with multiple transactions
//
func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	for value := 0; value < 10; value++ {
		t := strconv.FormatInt(time.Now().Unix(), 10)
		key := fmt.Sprintf("%s-%d", t, value)
		log.Infof("creating key(%s), value(%d)", key, value)
		if err := ctx.Create(key, fmt.Sprintf("%d", value), bluzelle.TestGasInfo(), nil); err != nil {
			log.Fatalf("%s", err)
		} else {
			log.Infof("created key")
		}
	}
}
