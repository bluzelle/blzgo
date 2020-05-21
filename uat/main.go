package main

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/vbstreetz/blzgo"
	"net/http"
	"os"
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

	switch request.Method {
	case "create":
		if len(request.Args) < 3 {
			abort(w, fmt.Errorf("key, value and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		value, ok := request.Args[1].(string)
		if !ok {
			abort(w, fmt.Errorf(VALUE_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[2])
		if err != nil {
			abort(w, err)
			return
		}
		var leaseInfo *bluzelle.LeaseInfo
		if len(request.Args) > 3 {
			if l, err := parseLeaseInfo(request.Args[3]); err != nil {
				abort(w, err)
				return
			} else {
				leaseInfo = l
			}
		}
		err = ctx.Create(key, value, gasInfo, leaseInfo)
		respond(&w, nil, err)
		return

	case "update":
		if len(request.Args) < 3 {
			abort(w, fmt.Errorf("key, value and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		value, ok := request.Args[1].(string)
		if !ok {
			abort(w, fmt.Errorf(VALUE_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[2])
		if err != nil {
			abort(w, err)
			return
		}
		var leaseInfo *bluzelle.LeaseInfo
		if len(request.Args) > 3 {
			if l, err := parseLeaseInfo(request.Args[3]); err != nil {
				abort(w, err)
				return
			} else {
				leaseInfo = l
			}
		}
		err = ctx.Update(key, value, gasInfo, leaseInfo)
		respond(&w, nil, err)
		return

	case "delete":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key and gas_info is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		err = ctx.Delete(key, gasInfo)
		respond(&w, nil, err)
		return

	case "rename":
		if len(request.Args) < 3 {
			abort(w, fmt.Errorf("key, new_key and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		newKey, ok := request.Args[1].(string)
		if !ok {
			abort(w, fmt.Errorf(NEW_KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[2])
		if err != nil {
			abort(w, err)
			return
		}
		err = ctx.Rename(key, newKey, gasInfo)
		respond(&w, nil, err)
		return

	case "delete_all":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("gas_info is required"))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[0])
		if err != nil {
			abort(w, err)
			return
		}
		err = ctx.DeleteAll(gasInfo)
		respond(&w, nil, err)
		return

	case "multi_update":
		kvsMap, ok := request.Args[0].([]interface{})
		if !ok {
			respond(&w, nil, fmt.Errorf("could not parse key values: %+v", request.Args[0]))
			return
		}

		keyValues := []*bluzelle.KeyValue{}

		for _, arg := range kvsMap {
			kv := &bluzelle.KeyValue{}
			kvMap, ok := arg.(map[string]interface{})
			if !ok {
				respond(&w, nil, fmt.Errorf("could not parse key values: %v", arg))
				return
			}
			keyArg := kvMap["key"]
			valueArg := kvMap["value"]

			if keyArg != nil {
				if key, ok := keyArg.(string); !ok {
					respond(&w, nil, fmt.Errorf("could not parse key in %v", kvMap))
					return
				} else {
					kv.Key = key
				}
			}
			if valueArg != nil {
				if value, ok := valueArg.(string); !ok {
					respond(&w, nil, fmt.Errorf("could not parse value in %v", kvMap))
					return
				} else {
					kv.Value = value
				}
			}

			keyValues = append(keyValues, kv)
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		err = ctx.MultiUpdate(keyValues, gasInfo)
		respond(&w, nil, err)
		return

	case "renew_lease":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		var leaseInfo *bluzelle.LeaseInfo
		if len(request.Args) > 2 {
			if l, err := parseLeaseInfo(request.Args[2]); err != nil {
				abort(w, err)
				return
			} else {
				leaseInfo = l
			}
		}
		err = ctx.RenewLease(key, gasInfo, leaseInfo)
		respond(&w, nil, err)
		return

	case "renew_lease_all":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("gas_info is required"))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[0])
		if err != nil {
			abort(w, err)
			return
		}
		var leaseInfo *bluzelle.LeaseInfo
		if len(request.Args) > 1 {
			if l, err := parseLeaseInfo(request.Args[1]); err != nil {
				abort(w, err)
				return
			} else {
				leaseInfo = l
			}
		}
		err = ctx.RenewLeaseAll(gasInfo, leaseInfo)
		respond(&w, nil, err)
		return

	//

	case "read":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		if len(request.Args) == 1 {
			v, err := ctx.Read(key)
			respond(&w, v, err)
			return
		} else {
			v, err := ctx.ProvenRead(key)
			respond(&w, v, err)
			return
		}
		return

	case "has":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		v, err := ctx.Has(key)
		respond(&w, v, err)
		return

	case "count":
		v, err := ctx.Count()
		respond(&w, v, err)
		return

	case "keys":
		v, err := ctx.Keys()
		respond(&w, v, err)
		return

	case "key_values":
		v, err := ctx.KeyValues()
		respond(&w, v, err)
		return

	case "get_lease":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		v, err := ctx.GetLease(key)
		respond(&w, v, err)
		return

	case "get_n_shortest_leases":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("n required"))
			return
		}
		n, ok := request.Args[0].(float64)
		if !ok {
			abort(w, fmt.Errorf(INVALID_VALUE_SPECIFIED))
			return
		}
		v, err := ctx.GetNShortestLeases(uint64(n))
		respond(&w, v, err)
		return

	//

	case "tx_read":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxRead(key, gasInfo)
		respond(&w, v, err)
		return

	case "tx_has":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxHas(key, gasInfo)
		respond(&w, v, err)
		return

	case "tx_count":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("gas_info required"))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[0])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxCount(gasInfo)
		respond(&w, v, err)
		return

	case "tx_keys":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("gas_info required"))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[0])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxKeys(gasInfo)
		respond(&w, v, err)
		return

	case "tx_key_values":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("gas_info required"))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[0])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxKeyValues(gasInfo)
		respond(&w, v, err)
		return

	case "tx_get_lease":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key and gas_info required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxGetLease(key, gasInfo)
		respond(&w, v, err)
		return

	case "tx_get_n_shortest_leases":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("N and gas_info required"))
			return
		}
		n, ok := request.Args[0].(float64)
		if !ok {
			abort(w, fmt.Errorf(INVALID_VALUE_SPECIFIED))
			return
		}
		gasInfo, err := parseGasInfo(request.Args[1])
		if err != nil {
			abort(w, err)
			return
		}
		v, err := ctx.TxGetNShortestLeases(uint64(n), gasInfo)
		respond(&w, v, err)
		return

	//

	case "account":
		v, err := ctx.Account()
		respond(&w, v, err)
		return

	case "version":
		v, err := ctx.Version()
		respond(&w, v, err)
		return

	//

	default:
		abort(w, fmt.Errorf("unsupported method: %s", request.Method))
		return
	}
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
	fmt.Fprint(*w, string(response))
	return
}
