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
const INVALID_LEASE_TIME = "Invalid lease time"
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
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and value are required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		value, ok := request.Args[1].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, leaseInfo, err := parseGasAndLeaseInfos(request, 2, 3)
		if err != nil {
			abort(w, err)
			return
		}
		if leaseInfo != nil && leaseInfo.ToBlocks() < 0 {
			abort(w, fmt.Errorf(INVALID_LEASE_TIME))
			return
		}
		if err := ctx.Create(key, value, gasInfo, leaseInfo); err != nil {
			abort(w, err)
		}
		return

	case "update":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and value are required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		value, ok := request.Args[1].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, leaseInfo, err := parseGasAndLeaseInfos(request, 2, 3)
		if err != nil {
			abort(w, err)
			return
		}
		if leaseInfo != nil && leaseInfo.ToBlocks() < 0 {
			abort(w, fmt.Errorf(INVALID_LEASE_TIME))
			return
		}
		if err := ctx.Update(key, value, gasInfo, leaseInfo); err != nil {
			abort(w, err)
		}
		return

	case "delete":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 1)
		if err != nil {
			abort(w, err)
			return
		}
		if err := ctx.Delete(key, gasInfo); err != nil {
			abort(w, err)
		}
		return

	case "rename":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and newkey are required"))
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
		gasInfo, _, err := parseGasAndLeaseInfos(request, 2, 3)
		if err != nil {
			abort(w, err)
			return
		}
		if err := ctx.Rename(key, newKey, gasInfo); err != nil {
			abort(w, err)
		}
		return

	case "delete_all":
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 0)
		if err != nil {
			abort(w, err)
			return
		}
		if err := ctx.DeleteAll(gasInfo); err != nil {
			abort(w, err)
		}
		return

	case "multi_update":
		keyValues := []*bluzelle.KeyValue{}
		for i := 0; i < ((len(request.Args) / 2) * 2); i = i + 2 {
			keyValues = append(keyValues, &bluzelle.KeyValue{
				Key:   request.Args[i].(string),
				Value: request.Args[i+1].(string),
			})
		}

		gasInfoIndex := len(request.Args)
		if gasInfoIndex%2 == 0 {
			gasInfoIndex = -1
		}

		if gasInfo, _, err := parseGasAndLeaseInfos(request, -1, gasInfoIndex); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.MultiUpdate(keyValues, gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	case "renew_lease":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and lease_info are required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, leaseInfo, err := parseGasAndLeaseInfos(request, 1, 2)
		if err != nil {
			abort(w, err)
			return
		}
		if leaseInfo != nil && leaseInfo.ToBlocks() < 0 {
			abort(w, fmt.Errorf(INVALID_LEASE_TIME))
			return
		}
		if err := ctx.RenewLease(key, gasInfo, leaseInfo); err != nil {
			abort(w, err)
		}
		return

	case "renew_all_leases":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and newkey are required"))
			return
		}
		gasInfo, leaseInfo, err := parseGasAndLeaseInfos(request, 0, 1)
		if err != nil {
			abort(w, err)
			return
		}
		if leaseInfo != nil && leaseInfo.ToBlocks() < 0 {
			abort(w, fmt.Errorf(INVALID_LEASE_TIME))
			return
		}
		if err := ctx.RenewAllLeases(gasInfo, leaseInfo); err != nil {
			abort(w, err)
		}
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
			if v, err := ctx.Read(key); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				if response, err := json.Marshal(v); err != nil {
					abort(w, fmt.Errorf("%s", err))
				} else {
					fmt.Fprintf(w, fmt.Sprintf("%s", response))
				}
			}
		} else {
			if v, err := ctx.ProvenRead(key); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				if response, err := json.Marshal(v); err != nil {
					abort(w, fmt.Errorf("%s", err))
				} else {
					fmt.Fprintf(w, fmt.Sprintf("%s", response))
				}
			}
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
		if v, err := ctx.Has(key); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "count":
		if v, err := ctx.Count(); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "keys":
		if v, err := ctx.Keys(); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "key_values":
		if v, err := ctx.KeyValues(); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
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
		if v, err := ctx.GetLease(key); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "get_n_shortest_leases":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("n is required"))
			return
		}
		n, ok := request.Args[0].(uint64)
		if !ok {
			abort(w, fmt.Errorf(INVALID_VALUE_SPECIFIED))
			return
		}
		if v, err := ctx.GetNShortestLeases(n); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	//

	case "tx_read":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 1)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxRead(key, gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "tx_has":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 1)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxHas(key, gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "tx_count":
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 0)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxCount(gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "tx_keys":
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 0)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxKeys(gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "tx_key_values":
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 0)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxKeyValues(gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "tx_get_lease":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key, ok := request.Args[0].(string)
		if !ok {
			abort(w, fmt.Errorf(KEY_MUST_BE_A_STRING))
			return
		}
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 1)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxGetLease(key, gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	case "tx_get_n_shortest_leases":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("n is required"))
			return
		}
		n, ok := request.Args[0].(uint64)
		if !ok {
			abort(w, fmt.Errorf(INVALID_VALUE_SPECIFIED))
			return
		}
		gasInfo, _, err := parseGasAndLeaseInfos(request, -1, 1)
		if err != nil {
			abort(w, err)
			return
		}
		if v, err := ctx.TxGetNShortestLeases(n, gasInfo); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			if response, err := json.Marshal(v); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%s", response))
			}
		}
		return

	//

	default:
		abort(w, fmt.Errorf("unsupported method: %s", request.Method))
		return
	}
}

func abort(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
}

func parseGasAndLeaseInfos(
	request *Request,
	gasInfoIndex int,
	leaseInfoIndex int,
) (
	*bluzelle.GasInfo,
	*bluzelle.LeaseInfo,
	error,
) {
	var gasInfo *bluzelle.GasInfo
	var leaseInfo *bluzelle.LeaseInfo

	if gasInfoIndex != -1 && len(request.Args) > gasInfoIndex {
		gasInfo = &bluzelle.GasInfo{}
		gasInfoArg := request.Args[gasInfoIndex]
		if gasInfoMap, ok := gasInfoArg.(map[string]interface{}); !ok {
			return nil, nil, fmt.Errorf("could not parse gasInfo: %v", gasInfoArg)
		} else {
			if gasInfoMap["max_gas"] != nil {
				if maxGas, ok := gasInfoMap["max_gas"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse gasInfo[maxGas]: %v", gasInfoMap["max_gas"])
				} else {
					gasInfo.MaxGas = int(maxGas)
				}
			}
			if gasInfoMap["max_fee"] != nil {
				if maxFee, ok := gasInfoMap["max_fee"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse gasInfo[maxFee]: %v", gasInfoMap["max_fee"])
				} else {
					gasInfo.MaxFee = int(maxFee)
				}
			}
			if gasInfoMap["gas_price"] != nil {
				if gasPrice, ok := gasInfoMap["gas_price"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse gasInfo[gasPrice]: %v", gasInfoMap["gas_price"])
				} else {
					gasInfo.GasPrice = int(gasPrice)
				}
			}
		}
	}

	if leaseInfoIndex != -1 && len(request.Args) > leaseInfoIndex {
		leaseInfo = &bluzelle.LeaseInfo{}
		leaseInfoArg := request.Args[leaseInfoIndex]
		if leaseInfoMap, ok := leaseInfoArg.(map[string]interface{}); !ok {
			return nil, nil, fmt.Errorf("could not parse leaseInfo: %v", leaseInfoArg)
		} else {
			if leaseInfoMap["days"] != nil {
				if days, ok := leaseInfoMap["days"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse leaseInfo[days]: %v", leaseInfoMap["days"])
				} else {
					leaseInfo.Days = int64(days)
				}
			}
			if leaseInfoMap["hours"] != nil {
				if hours, ok := leaseInfoMap["hours"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse leaseInfo[hours]: %v", leaseInfoMap["hours"])
				} else {
					leaseInfo.Hours = int64(hours)
				}
			}
			if leaseInfoMap["minutes"] != nil {
				if minutes, ok := leaseInfoMap["minutes"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse leaseInfo[minutes]: %v", leaseInfoMap["minutes"])
				} else {
					leaseInfo.Minutes = int64(minutes)
				}
			}
			if leaseInfoMap["seconds"] != nil {
				if seconds, ok := leaseInfoMap["seconds"].(float64); !ok {
					return nil, nil, fmt.Errorf("could not parse leaseInfo[seconds]: %v", leaseInfoMap["seconds"])
				} else {
					leaseInfo.Seconds = int64(seconds)
				}
			}
		}
	}

	return gasInfo, leaseInfo, nil
}
