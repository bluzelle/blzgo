package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

type Result struct {
	Error error
	Key   string
}

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

	no_of_creates := 5
	creates := make(chan Result, no_of_creates)

	for value := 0; value < no_of_creates; value++ {
		t := strconv.FormatInt(time.Now().Unix(), 10)
		key := fmt.Sprintf("%s-%d", t, value)
		log.Infof("creating key(%s), value(%d)", key, value)
		go create(ctx, key, fmt.Sprintf("%d", value), creates)
	}

	i := 0
	for ret := range creates {
		if ret.Error != nil {
			log.Fatalf("error(%s): %s", ret.Key, err)
		} else {
			log.Infof("created key(%s)", ret.Key)
		}
		i++
		if i == no_of_creates {
			close(creates)
		}
	}
}

func create(ctx *bluzelle.Client, key string, value string, c chan Result) {
	err := ctx.Create(key, value, bluzelle.GetTestGasInfo(), nil)
	c <- Result{
		err, key,
	}
}
