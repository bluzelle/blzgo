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
		if leaseInfo, gasInfo, err := parseLeaseAndGasInfo(request, 2, 3); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.Create(request.Args[0].(string), request.Args[1].(string), leaseInfo, gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	case "update":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and value are required"))
			return
		}
		if leaseInfo, gasInfo, err := parseLeaseAndGasInfo(request, 2, 3); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.Update(request.Args[0].(string), request.Args[1].(string), leaseInfo, gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	case "delete":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 1); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.Delete(request.Args[0].(string), gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	case "rename":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and newkey are required"))
			return
		}
		if _, gasInfo, err := parseLeaseAndGasInfo(request, 2, 3); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.Rename(request.Args[0].(string), request.Args[1].(string), gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	case "delete_all":
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 0); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.DeleteAll(gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

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

		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, gasInfoIndex); err != nil {
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
		if leaseInfo, gasInfo, err := parseLeaseAndGasInfo(request, 1, 2); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.RenewLease(request.Args[0].(string), leaseInfo, gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	case "renew_all_leases":
		if len(request.Args) < 2 {
			abort(w, fmt.Errorf("both key and newkey are required"))
			return
		}
		if leaseInfo, gasInfo, err := parseLeaseAndGasInfo(request, 0, 1); err != nil {
			abort(w, err)
			return
		} else {
			if err := ctx.RenewAllLeases(leaseInfo, gasInfo); err != nil {
				abort(w, err)
			}
			return
		}

	//

	case "read":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		key := request.Args[0].(string)
		if len(request.Args) == 1 {
			if v, err := ctx.Read(key); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		} else {
			if v, err := ctx.ProvenRead(key); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "has":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		if v, err := ctx.Has(request.Args[0].(string)); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}
		return

	case "count":
		if v, err := ctx.Count(); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}
		return

	case "keys":
		if v, err := ctx.Keys(); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}
		return

	case "key_values":
		if v, err := ctx.KeyValues(); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}
		return

	case "get_lease":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		if v, err := ctx.GetLease(request.Args[0].(string)); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}
		return

	case "get_n_shortest_leases":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("n is required"))
			return
		}
		if v, err := ctx.GetNShortestLeases(request.Args[0].(uint64)); err != nil {
			abort(w, fmt.Errorf("%s", err))
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}
		return

	//

	case "tx_read":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}

		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 1); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxRead(request.Args[0].(string), gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "tx_has":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}

		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 1); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxHas(request.Args[0].(string), gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "tx_count":
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 0); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxCount(gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "tx_keys":
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 0); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxKeys(gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "tx_key_values":
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 0); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxKeyValues(gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "tx_get_lease":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("key is required"))
			return
		}
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 1); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxGetLease(request.Args[0].(string), gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
			}
		}
		return

	case "tx_get_n_shortest_leases":
		if len(request.Args) < 1 {
			abort(w, fmt.Errorf("n is required"))
			return
		}
		if _, gasInfo, err := parseLeaseAndGasInfo(request, -1, 1); err != nil {
			abort(w, err)
			return
		} else {
			if v, err := ctx.TxGetNShortestLeases(request.Args[0].(uint64), gasInfo); err != nil {
				abort(w, fmt.Errorf("%s", err))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
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

func parseLeaseAndGasInfo(
	request *Request,
	gasInfoIndex int,
	leaseInfoIndex int,
) (
	*bluzelle.LeaseInfo,
	*bluzelle.GasInfo,
	error,
) {
	var leaseInfo *bluzelle.LeaseInfo
	var gasInfo *bluzelle.GasInfo

	if leaseInfoIndex != -1 && len(request.Args) > leaseInfoIndex {
		leaseInfo := &bluzelle.LeaseInfo{}
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

		if gasInfoIndex != -1 && len(request.Args) > gasInfoIndex {
			gasInfo := &bluzelle.GasInfo{}
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
	}

	return leaseInfo, gasInfo, nil
}
