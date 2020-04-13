package bluzelle

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const TX_COMMAND string = "/txs"
const TOKEN_NAME string = "ubnt"

//
// JSON struct keys are ordered alphabetically
//
type TransactionInitRequestBaseReq struct {
	From    string `json:"from"`
	ChainId string `json:"chain_id"`
}

type KeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type TransactionInitRequest struct {
	BaseReq   *TransactionInitRequestBaseReq `json:"BaseReq"`
	UUID      string                         `json:"UUID"`
	Key       string                         `json:"Key,omitempty"`
	KeyValues []*KeyValue                    `json:"KeyValues,omitempty"`
	NewKey    string                         `json:"NewKey,omitempty"`
	Value     string                         `json:"Value,omitempty"`
	Owner     string                         `json:"Owner"`
}

type TransactionInitResponseValueMsgValue struct {
	Key       string      `json:"Key,omitempty"`
	KeyValues []*KeyValue `json:"KeyValues,omitempty"`
	NewKey    string      `json:"NewKey,omitempty"`
	Owner     string      `json:"Owner"`
	UUID      string      `json:"UUID"`
	Value     string      `json:"Value,omitempty"`
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

type GasInfo struct {
	MaxGas   int
	MaxFee   int
	GasPrice int
}

type Transaction struct {
	Key                string
	Value              string
	KeyValues          []*KeyValue
	NewKey             string
	ApiRequestMethod   string
	ApiRequestEndpoint string
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
	if err == nil {
		transaction.Client.account.Sequence += 1
	}
	transaction.Done(b, err)
}

func (transaction *Transaction) Init() (*TransactionInitResponseValue, error) {
	req := &TransactionInitRequest{
		BaseReq: &TransactionInitRequestBaseReq{
			From:    transaction.Client.options.Address,
			ChainId: transaction.Client.options.ChainId,
		},
		UUID:      transaction.Client.options.UUID,
		Key:       transaction.Key,
		KeyValues: transaction.KeyValues,
		NewKey:    transaction.NewKey,
		Owner:     transaction.Client.options.Address,
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

	transaction.Client.Infof("txn init %+v", string(body))

	res := &TransactionInitResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func (transaction *Transaction) Broadcast(data *TransactionInitResponseValue) ([]byte, error) {
	// Check the gas info
	feeGas, err := strconv.Atoi(data.Fee.Gas)
	if err != nil {
		transaction.Client.Errorf("failed to pass gas to int(%)", data.Fee.Gas)
	}
	gasInfo := transaction.Client.options.GasInfo
	if gasInfo.MaxGas != 0 && feeGas > gasInfo.MaxGas {
		data.Fee.Gas = strconv.Itoa(gasInfo.MaxGas)
	}
	if gasInfo.MaxFee != 0 {
		data.Fee.Amount = []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(gasInfo.MaxFee)},
		}
	} else if gasInfo.GasPrice != 0 {
		data.Fee.Amount = []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(feeGas * gasInfo.GasPrice)},
		}
	}

	// Create transaction payload
	txn := &TransactionBroadcastRequestTransaction{
		Msg:  data.Msg,
		Fee:  data.Fee,
		Memo: makeRandomString(32),
	}

	// Sign transaction
	sig, err := transaction.Sign(txn)
	if err != nil {
		return nil, err
	}
	txn.Signatures = []*TransactionSignature{sig}
	txn.Signature = sig

	// Broadcast transaction
	req := &TransactionBroadcastRequest{
		Transaction: txn,
		Mode:        "block",
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn broadcast request %+v", string(reqBytes))
	body, err := transaction.Client.APIMutate("POST", TX_COMMAND, reqBytes)
	if err != nil {
		return nil, err
	}

	// Read transaction broadcast response
	res := &TransactionBroadcastResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn broadcast response %+v", res)
	if res.Code != 0 {
		return nil, fmt.Errorf("%s", res.RawLog)
	}
	decodedData, err := hex.DecodeString(res.Data)
	return decodedData, err
}

func (transaction *Transaction) Sign(req *TransactionBroadcastRequestTransaction) (*TransactionSignature, error) {
	// pubKeyValue
	pubKeyValue := base64.StdEncoding.EncodeToString(transaction.Client.privateKey.PubKey().SerializeCompressed())

	// accountNumber
	accountNumber := strconv.Itoa(transaction.Client.account.AccountNumber)

	// Sequence
	seq := strconv.Itoa(transaction.Client.account.Sequence)

	// Calculate the SHA256 of the payload object
	payload := &TransactionBroadcastSignPayload{
		AccountNumber: strconv.Itoa(transaction.Client.account.AccountNumber),
		ChainId:       transaction.Client.options.ChainId,
		Memo:          req.Memo,
		Sequence:      seq,
		Fee:           req.Fee, // already sorted by key
		Msgs:          req.Msg, // already sorted by key
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn sign %+v", string(payloadBytes))
	sig := ""
	hash := tmcrypto.Sha256(payloadBytes)
	if s, err := transaction.Client.privateKey.Sign(hash); err != nil {
		return nil, err
	} else {
		sig = base64.StdEncoding.EncodeToString(serializeSig(s))
	}

	return &TransactionSignature{
		PubKey: &TransactionSignaturePubKey{
			Type:  tmsecp256k1.PubKeyAminoName,
			Value: pubKeyValue,
		},
		Signature:     sig,
		AccountNumber: accountNumber,
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

// https://github.com/tendermint/tendermint/blob/ef56e6661121a7f8d054868689707cd817f27b24/crypto/secp256k1/secp256k1_nocgo.go#L60-L70
// Serialize signature to R || S.
// R, S are padded to 32 bytes respectively.
func serializeSig(sig *btcec.Signature) []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}
