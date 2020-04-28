package bluzelle

import (
	"github.com/apex/log"
	clih "github.com/apex/log/handlers/cli"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Test struct {
	Client *Client
	Key1   string
	Key2   string
	Key3   string
	Value1 string
	Value2 string
	Value3 string
}

func (ctx *Test) SetUp() error {
	SetupLogging()
	LoadEnv()

	c, err := NewTestClient()
	if err != nil {
		return err
	} else {
		ctx.Client = c
	}

	ctx.Key1 = strconv.FormatInt(100+time.Now().Unix(), 10)
	ctx.Key2 = strconv.FormatInt(200+time.Now().Unix(), 10)
	ctx.Key3 = strconv.FormatInt(300+time.Now().Unix(), 10)

	ctx.Value1 = "foo"
	ctx.Value2 = "bar"
	ctx.Value3 = "baz"

	return nil
}

func (ctx *Test) TearDown() error {
	return nil
}

func NewTestClient() (*Client, error) {
	debug := false
	if d, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
		debug = d
	}

	// create client
	options := &Options{
		Address:  os.Getenv("ADDRESS"),
		Mnemonic: os.Getenv("MNEMONIC"),
		UUID:     os.Getenv("UUID"),
		Endpoint: os.Getenv("ENDPOINT"),
		ChainId:  os.Getenv("CHAIN_ID"),
		GasInfo: &GasInfo{
			MaxFee: 4000001,
		},
		Debug: debug,
	}
	ctx, err := NewClient(options)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func SetupLogging() {
	log.SetHandler(clih.Default)
	log.SetLevel(log.DebugLevel)
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		// log.Errorf("%s", err)
		if err := godotenv.Load("../.env"); err != nil { // when running tests
			// log.Errorf("%s", err)
		}
	}
}
