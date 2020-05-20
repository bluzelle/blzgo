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

func (ctx *Test) TestSetUp() error {
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

func (ctx *Test) TestTearDown() error {
	return nil
}

func NewTestClient() (*Client, error) {
	debug := false
	if d, err := strconv.ParseBool(os.Getenv("DEBUG")); err == nil {
		debug = d
	}

	// create client
	options := &Options{
		Mnemonic: os.Getenv("MNEMONIC"),
		UUID:     os.Getenv("UUID"),
		Endpoint: os.Getenv("ENDPOINT"),
		ChainId:  os.Getenv("CHAIN_ID"),
		Debug:    debug,
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

func TestGasInfo() *GasInfo {
	return &GasInfo{
		MaxFee: 4000001,
	}
}

func TestAddress() string {
	return os.Getenv("ADDRESS")
}

func sanitizeString(s string) string {
	return s
	// return re.ReplaceAllStringFunc(s, sanitizeStringToken)
	// s.gsub(/([&<>])/) { |token|
	// 	"\\u00#{token[0].ord.to_s(16)}"
	// }
}

// func sanitizeStringToken(token string) {
// 	return "\\u00#{token[0].ord.to_s(16)}"
// }

func encodeSafe(s string) string {
	return s
	// return re.ReplaceAllStringFunc(s, sanitizeStringToken)
}

// func encodeSafeToken(token string) {
// 	return "%#{token[0].ord.to_s(16)}"
// }
