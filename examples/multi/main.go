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

	creates := make(chan error)

	for value := 0; value < 10; value++ {
		t := strconv.FormatInt(time.Now().Unix(), 10)
		key := fmt.Sprintf("%s-%d", t, value)
		log.Infof("creating key(%s), value(%d)", key, value)
		go create(ctx, key, fmt.Sprintf("%d", value), creates)
	}

	for err := range creates {
		if err != nil {
			log.Fatalf("%s", err)
		} else {
			log.Infof("created key")
		}
	}
}

func create(ctx *bluzelle.Client, key string, value string, c chan error) {
	err := ctx.Create(key, value, bluzelle.GetTestGasInfo(), nil)
	c <- err
}
