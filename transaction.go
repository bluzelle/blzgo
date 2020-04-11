package bluzelle

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const TX_COMMAND string = "/txs"
const TOKEN_NAME string = "ubnt"

//

type TransactionInitRequestBaseReq struct {
	From    string `json:"from"`
	ChainId string `json:"chain_id"`
}

type TransactionInitRequestKeyValue struct {
	Key       string                            `json:"key,omitempty"`
	KeyValues []*TransactionInitRequestKeyValue `json:"KeyValues,omitempty"`
	NewKey    string                            `json:"NewKey,omitempty"`
	Value     string                            `json:"value,omitempty"`
}

type TransactionInitRequest struct {
	BaseReq   *TransactionInitRequestBaseReq    `json:"BaseReq"`
	UUID      string                            `json:"UUID"`
	Key       string                            `json:"Key,omitempty"`
	KeyValues []*TransactionInitRequestKeyValue `json:"KeyValues,omitempty"`
	NewKey    string                            `json:"NewKey,omitempty"`
	Value     string                            `json:"Value,omitempty"`
	Owner     string                            `json:"Owner"`
}

type TransactionInitResponseValueMsgValue struct {
	Key       string                            `json:"Key,omitempty"`
	KeyValues []*TransactionInitRequestKeyValue `json:"KeyValues,omitempty"`
	NewKey    string                            `json:"NewKey,omitempty"`
	Owner     string                            `json:"Owner"`
	UUID      string                            `json:"UUID"`
	Value     string                            `json:"Value,omitempty"`
}

type TransactionInitResponseValueMsg struct {
	Type  string                                `json:"type"`
	Value *TransactionInitResponseValueMsgValue `json:"value"`
}

type TransactionFeeAmount struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type TransactionFee struct {
	Amount []*TransactionFeeAmount `json:"amount"`
	Gas    string                  `json:"gas"`
}

type TransactionSignaturePubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TransactionSignature struct {
	AccountNumber string                      `json:"account_number"`
	PubKey        *TransactionSignaturePubKey `json:"pub_key"`
	Sequence      string                      `json:"sequence"`
	Signature     string                      `json:"signature"`
}

type TransactionInitResponseValue struct {
	Fee        *TransactionFee                    `json:"fee"`
	Memo       string                             `json:"memo"`
	Msg        []*TransactionInitResponseValueMsg `json:"msg"`
	Signatures []*TransactionSignature            `json:"signatures"`
}

type TransactionInitResponse struct {
	Type  string                        `json:"type"`
	Value *TransactionInitResponseValue `json:"value"`
}

//

type TransactionBroadcastSignPayload struct {
	AccountNumber string                             `json:"account_number"`
	ChainId       string                             `json:"chain_id"`
	Fee           *TransactionFee                    `json:"fee"`
	Memo          string                             `json:"memo"`
	Msgs          []*TransactionInitResponseValueMsg `json:"msgs"`
	Sequence      string                             `json:"sequence"`
}

type TransactionBroadcastRequestTransaction struct {
	Fee        *TransactionFee                    `json:"fee"`
	Memo       string                             `json:"memo"`
	Msg        []*TransactionInitResponseValueMsg `json:"msg"`
	Signature  *TransactionSignature              `json:"signature"`
	Signatures []*TransactionSignature            `json:"signatures"`
}

type TransactionBroadcastRequest struct {
	Transaction *TransactionBroadcastRequestTransaction `json:"tx"`
	Mode        string                                  `json:"mode"`
}

type TransactionBroadcastResponse struct {
	Height    string `json:"height"`
	TxHash    string `json:"txhash"`
	Data      string `json:"data"`
	Codespace string `json:"codespace"`
	Code      int    `json:"code"`
	RawLog    string `json:"raw_log"`
	GasWanted string `json:"gas_wanted"`
}

//

type Transaction struct {
	Key                string
	Value              string
	KeyValues          []*TransactionInitRequestKeyValue
	NewKey             string
	Address            string
	UUID               string
	ApiRequestMethod   string
	ApiRequestEndpoint string
	GasInfo            *GasInfo
	ChainId            string
	Client             *Client

	done   chan bool
	result []byte
	err    error
}

func (transaction *Transaction) Done(result []byte, err error) {
	transaction.result = result
	transaction.err = err
	transaction.done <- true
	close(transaction.done)
}

func (transaction *Transaction) Send() {
	res, err := transaction.Init()
	if err != nil {
		transaction.Done(nil, err)
		return
	}
	b, err := transaction.Broadcast(res)
	transaction.Done(b, err)
}

func (transaction *Transaction) Init() (*TransactionInitResponseValue, error) {
	req := &TransactionInitRequest{
		BaseReq: &TransactionInitRequestBaseReq{
			From:    transaction.Address,
			ChainId: transaction.ChainId,
		},
		UUID:      transaction.UUID,
		Key:       transaction.Key,
		KeyValues: transaction.KeyValues,
		NewKey:    transaction.NewKey,
		Owner:     transaction.Address,
		Value:     transaction.Value,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn init %+v", string(reqBytes))
	body, err := transaction.Client.APIMutate(transaction.ApiRequestMethod, transaction.ApiRequestEndpoint, reqBytes)
	if err != nil {
		return nil, err
	}

	res := &TransactionInitResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func (transaction *Transaction) Broadcast(data *TransactionInitResponseValue) ([]byte, error) {
	// set the gas info
	feeGas, err := strconv.Atoi(data.Fee.Gas)
	if err != nil {
		transaction.Client.Errorf("failed to pass gas to int(%)", data.Fee.Gas)
	}

	if transaction.GasInfo.MaxGas != 0 && feeGas > transaction.GasInfo.MaxGas {
		data.Fee.Gas = strconv.Itoa(transaction.GasInfo.MaxGas)
	}

	if transaction.GasInfo.MaxFee != 0 {
		data.Fee.Amount = []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(transaction.GasInfo.MaxFee)},
		}
	} else if transaction.GasInfo.GasPrice != 0 {
		data.Fee.Amount = []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(feeGas * transaction.GasInfo.GasPrice)},
		}
	}

	// broadcast
	txn := &TransactionBroadcastRequestTransaction{
		Msg:  data.Msg,
		Fee:  data.Fee,
		Memo: makeRandomString(32),
	}

	// sign
	sig, err := transaction.Sign(txn)
	if err != nil {
		return nil, err
	}
	txn.Signatures = []*TransactionSignature{sig}
	txn.Signature = sig

	req := &TransactionBroadcastRequest{
		Transaction: txn,
		Mode:        "block",
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn broadcast %+v", string(reqBytes))
	body, err := transaction.Client.APIMutate("POST", TX_COMMAND, reqBytes)
	if err != nil {
		return nil, err
	}

	res := &TransactionBroadcastResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	transaction.Client.Infof("broadcast %+v", res)

	if res.Code != 0 {
		return nil, fmt.Errorf("%s", res.RawLog)
	}
	decodedData, err := hex.DecodeString(res.Data)
	return decodedData, err
}

func (transaction *Transaction) Sign(req *TransactionBroadcastRequestTransaction) (*TransactionSignature, error) {
	pubKeyString := ""
	pubKey := transaction.Client.PrivateKey.PubKey()
	if b, err := hex.DecodeString(fmt.Sprintf("%x", secp256k1.CompressPubkey(pubKey.X, pubKey.Y))); err != nil {
		return nil, err
	} else {
		pubKeyString = base64.StdEncoding.EncodeToString(b)
	}

	seq := strconv.Itoa(transaction.Client.Account.Sequence) // + 1

	// Calculate the SHA256 of the payload object
	payload := &TransactionBroadcastSignPayload{
		AccountNumber: strconv.Itoa(transaction.Client.Account.AccountNumber),
		ChainId:       transaction.Client.Options.ChainId,
		Memo:          req.Memo,
		Sequence:      seq,
		Fee:           req.Fee, // alreayd sorted
		Msgs:          req.Msg, // alreayd sorted
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn sign %+v", string(payloadBytes))
	sig := ""
	hash := hashSha256(payloadBytes)
	if s, err := transaction.Client.PrivateKey.Sign(hash); err != nil {
		return nil, err
	} else {
		// We have to convert the signature to the format that Tendermint uses
		g := []byte{}
		g = append(s.R.Bytes(), s.S.Bytes()...)
		sig = base64.StdEncoding.EncodeToString(g)
	}

	// transaction.Client.Infof("hash %x", hash)
	// transaction.Client.Infof("sig %s", sig)

	return &TransactionSignature{
		PubKey: &TransactionSignaturePubKey{
			Type:  "tendermint/PubKeySecp256k1",
			Value: pubKeyString,
		},
		Signature:     sig,
		AccountNumber: strconv.Itoa(transaction.Client.Account.AccountNumber),
		Sequence:      seq,
	}, nil
}

func makeRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
