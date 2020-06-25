package main

import (
	"os"
	"strconv"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("both key and value are required")
	}

	ctx, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	}

	key := args[0]
	value := args[1]
	leaseInfo := &bluzelle.LeaseInfo{}

	if len(args) > 2 {
		if lease, err := strconv.Atoi(args[2]); err != nil {
			log.Fatalf("could not parse provided lease(%s)", args[2])
		} else {
			leaseInfo.Seconds = int64(lease)
		}
	}

	log.Infof("creating key(%s), val(%s), lease(%ds)...", key, value, leaseInfo.ToBlocks())

	if err := ctx.Create(key, value, bluzelle.TestGasInfo(), leaseInfo); err != nil {
		log.Fatalf("%s", err)
	} else {
		log.Infof("created key")
	}
}
