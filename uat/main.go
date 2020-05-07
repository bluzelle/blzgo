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
	Method string   `json:"method"`
	Args   []string `json:"args"`
}

type LeaseInfo struct {
	Days    int64 `json:"days"`
	Hours   int64 `json:"hours"`
	Minutes int64 `json:"minutes"`
	Seconds int64 `json:"seconds"`
}

type GasInfo struct {
	MaxGas   int `json:"max_gas"`
	MaxFee   int `json:"max_fee"`
	GasPrice int `json:"gas_price"`
}

func uat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method.", 405)
		return
	}

	var request *Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case "read":
		if len(request.Args) != 1 {
			http.Error(w, fmt.Sprintf("invalid number of keys passed: %+v", request.Args), 500)
			return
		}
		if v, err := ctx.Read(request.Args[0]); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), 500)
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%v\n", v))
		}

	case "create":
		if len(request.Args) < 2 {
			http.Error(w, fmt.Sprintf("both key and value are required"), 500)
			return
		}
		if err := ctx.Create(request.Args[0], request.Args[1], nil, nil); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), 500)
		}

	default:
		http.Error(w, fmt.Sprintf("unsupported method: %s", request.Method), 500)
	}
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
