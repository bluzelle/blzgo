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
		log.Fatalf("both key and lease are required")
	}

	ctx, err := util.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	lease := 0
	if l, err := strconv.Atoi(args[1]); err != nil {
		log.Fatalf("could not parse provided lease(%s)", args[1])
	} else {
		lease = l
	}

	log.Infof("renewing key(%s), lease(%d)...", key, lease)

	if err := ctx.RenewLease(key, int64(lease)); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("renewed lease")
	}
}
