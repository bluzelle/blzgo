package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/apex/log"
	bluzelle "github.com/bluzelle/blzgo"
)

var ctx *bluzelle.Client

const KEY_MUST_BE_A_STRING = "Key must be a string"
const NEW_KEY_MUST_BE_A_STRING = "New key must be a string"
const VALUE_MUST_BE_A_STRING = "Value must be a string"
const ALL_KEYS_MUST_BE_STRINGS = "All keys must be strings"
const ALL_VALUES_MUST_BE_STRINGS = "All values must be strings"

// const INVALID_LEASE_TIME = "Invalid lease time"
const INVALID_VALUE_SPECIFIED = "Invalid value specified"
const ADDRESS_MUST_BE_A_STRING = "address must be a string"
const MNEMONIC_MUST_BE_A_STRING = "mnemonic must be a string"
const UUID_MUST_BE_A_STRING = "uuid must be a string"
const INVALID_TRANSACTION = "Invalid transaction."

type Request struct {
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
}

func main() {
	bluzelle.SetupLogging()
	bluzelle.LoadEnv()

	c, err := bluzelle.NewTestClient()
	if err != nil {
		log.Fatalf("%s", err)
	} else {
		ctx = c
	}

	http.HandleFunc("/", uat)
	port := os.Getenv("PORT")
	if port == "" {
		port = "4562"
	}
	log.Infof("serving at :%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func uat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		abort(w, fmt.Errorf("invalid request method."))
		return
	}

	var request *Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		abort(w, err)
		return
	}

	inputs := make([]reflect.Value, len(request.Args))
	for i, _ := range request.Args {
		inputs[i] = reflect.ValueOf(request.Args[i])
	}

	switch request.Method {
	case "Create":
		if v, err := parseGasInfo(request.Args[2]); err != nil {
			abort(w, err)
			return
		} else if v != nil {
			inputs[2] = reflect.ValueOf(v)
		}
		if v, err := parseLeaseInfo(request.Args[3]); err != nil {
			abort(w, err)
			return
		} else if v != nil {
			inputs[3] = reflect.ValueOf(v)
		}
	default:
		abort(w, fmt.Errorf("unsupported method: %s", request.Method))
		return
	}

	result := reflect.ValueOf(ctx).MethodByName(request.Method).Call(inputs)
	errIndex := len(result)
	if e := result[errIndex-1].Interface(); e != nil {
		err = e.(error)
	} else {
		err = nil
	}
	respond(&w, result[0].Interface(), err)
}

func parseGasInfo(arg interface{}) (*bluzelle.GasInfo, error) {
	gasInfo := &bluzelle.GasInfo{}

	gasInfoMap, ok := arg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse gasInfo: %v", arg)
	}

	if gasInfoMap["max_gas"] != nil {
		if maxGas, ok := gasInfoMap["max_gas"].(float64); !ok {
			return nil, fmt.Errorf("could not parse gasInfo[maxGas]: %v", gasInfoMap["max_gas"])
		} else {
			gasInfo.MaxGas = int(maxGas)
		}
	}
	if gasInfoMap["max_fee"] != nil {
		if maxFee, ok := gasInfoMap["max_fee"].(float64); !ok {
			return nil, fmt.Errorf("could not parse gasInfo[maxFee]: %v", gasInfoMap["max_fee"])
		} else {
			gasInfo.MaxFee = int(maxFee)
		}
	}
	if gasInfoMap["gas_price"] != nil {
		if gasPrice, ok := gasInfoMap["gas_price"].(float64); !ok {
			return nil, fmt.Errorf("could not parse gasInfo[gasPrice]: %v", gasInfoMap["gas_price"])
		} else {
			gasInfo.GasPrice = int(gasPrice)
		}
	}

	return gasInfo, nil
}

func parseLeaseInfo(arg interface{}) (*bluzelle.LeaseInfo, error) {
	leaseInfo := &bluzelle.LeaseInfo{}

	leaseInfoMap, ok := arg.(map[string]interface{})
	if !ok {
		// return nil, fmt.Errorf("could not parse leaseInfo: %v", arg)
		return leaseInfo, nil
	}

	if leaseInfoMap["days"] != nil {
		if days, ok := leaseInfoMap["days"].(float64); !ok {
			return nil, fmt.Errorf("could not parse leaseInfo[days]: %v", leaseInfoMap["days"])
		} else {
			leaseInfo.Days = int64(days)
		}
	}
	if leaseInfoMap["hours"] != nil {
		if hours, ok := leaseInfoMap["hours"].(float64); !ok {
			return nil, fmt.Errorf("could not parse leaseInfo[hours]: %v", leaseInfoMap["hours"])
		} else {
			leaseInfo.Hours = int64(hours)
		}
	}
	if leaseInfoMap["minutes"] != nil {
		if minutes, ok := leaseInfoMap["minutes"].(float64); !ok {
			return nil, fmt.Errorf("could not parse leaseInfo[minutes]: %v", leaseInfoMap["minutes"])
		} else {
			leaseInfo.Minutes = int64(minutes)
		}
	}
	if leaseInfoMap["seconds"] != nil {
		if seconds, ok := leaseInfoMap["seconds"].(float64); !ok {
			return nil, fmt.Errorf("could not parse leaseInfo[seconds]: %v", leaseInfoMap["seconds"])
		} else {
			leaseInfo.Seconds = int64(seconds)
		}
	}

	return leaseInfo, nil
}

func abort(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
}

func respond(w *http.ResponseWriter, v interface{}, err error) {
	if err != nil {
		abort(*w, fmt.Errorf("%s", err))
		return
	}
	response, err := json.Marshal(v)
	if err != nil {
		abort(*w, fmt.Errorf("%s", err))
	}
	fmt.Fprintf(*w, fmt.Sprintf("%s", response))
	return
}
