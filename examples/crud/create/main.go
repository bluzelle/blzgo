package main

import (
	"github.com/apex/log"
	util "github.com/vbstreetz/blzgo"
	"os"
	"strconv"
)

func main() {
	util.SetupLogging()
	util.LoadEnv()

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("both key and value are required")
	}

	ctx, err := util.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	value := args[1]
	lease := 0

	if len(args) > 2 {
		if l, err := strconv.Atoi(args[2]); err != nil {
			log.Fatalf("could not parse provided lease(%s)", args[2])
		} else {
			lease = l
		}
	}

	log.Infof("creating key(%s), val(%s), lease(%d)...", key, value, lease)

	if err := ctx.Create(key, value, int64(lease)); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}
}
