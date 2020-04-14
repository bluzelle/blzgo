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
const BROADCAST_MAX_RETRIES = 10
const BROADCAST_RETRY_INTERVAL time.Duration = time.Second

//
// JSON struct keys are ordered alphabetically
//

type KeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type GasInfo struct {
	MaxGas   int
	MaxFee   int
	GasPrice int
}

//

type TransactionFeeAmount struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type TransactionFee struct {
	Amount []*TransactionFeeAmount `json:"amount"`
	Gas    string                  `json:"gas"`
}

//

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

//

type TransactionMsgValue struct {
	Key       string      `json:"Key,omitempty"`
	KeyValues []*KeyValue `json:"KeyValues,omitempty"`
	NewKey    string      `json:"NewKey,omitempty"`
	Owner     string      `json:"Owner"`
	UUID      string      `json:"UUID"`
	Value     string      `json:"Value,omitempty"`
}

type TransactionMsg struct {
	Type  string               `json:"type"`
	Value *TransactionMsgValue `json:"value"`
}

//

type TransactionBroadcastPayload struct {
	Fee        *TransactionFee         `json:"fee"`
	Memo       string                  `json:"memo"`
	Msg        []*TransactionMsg       `json:"msg"`
	Signatures []*TransactionSignature `json:"signatures"`
}

//

type TransactionValidateRequestBaseReq struct {
	From    string `json:"from"`
	ChainId string `json:"chain_id"`
}

type TransactionValidateRequest struct {
	BaseReq   *TransactionValidateRequestBaseReq `json:"BaseReq"`
	UUID      string                             `json:"UUID"`
	Key       string                             `json:"Key,omitempty"`
	KeyValues []*KeyValue                        `json:"KeyValues,omitempty"`
	NewKey    string                             `json:"NewKey,omitempty"`
	Value     string                             `json:"Value,omitempty"`
	Owner     string                             `json:"Owner"`
}

type TransactionValidateResponse struct {
	Type  string                       `json:"type"`
	Value *TransactionBroadcastPayload `json:"value"`
}

//

type TransactionBroadcastRequest struct {
	Transaction *TransactionBroadcastPayload `json:"tx"`
	Mode        string                       `json:"mode"`
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

type TransactionBroadcastPayloadSignPayload struct {
	AccountNumber string            `json:"account_number"`
	ChainId       string            `json:"chain_id"`
	Fee           *TransactionFee   `json:"fee"`
	Memo          string            `json:"memo"`
	Msgs          []*TransactionMsg `json:"msgs"`
	Sequence      string            `json:"sequence"`
}

//

type Transaction struct {
	Key                string
	Value              string
	KeyValues          []*KeyValue
	NewKey             string
	ApiRequestMethod   string
	ApiRequestEndpoint string
	Client             *Client

	broadcastRetries int
}

// Get required min gas
func (transaction *Transaction) Validate() (*TransactionBroadcastPayload, error) {
	req := &TransactionValidateRequest{
		BaseReq: &TransactionValidateRequestBaseReq{
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

	res := &TransactionValidateResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func (transaction *Transaction) Broadcast(txn *TransactionBroadcastPayload) ([]byte, error) {
	// Set memo
	txn.Memo = makeRandomString(32)

	// Set fee
	feeGas, err := strconv.Atoi(txn.Fee.Gas)
	if err != nil {
		transaction.Client.Errorf("failed to pass gas to int(%)", txn.Fee.Gas)
	}
	gasInfo := transaction.Client.options.GasInfo
	if gasInfo.MaxGas != 0 && feeGas > gasInfo.MaxGas {
		txn.Fee.Gas = strconv.Itoa(gasInfo.MaxGas)
	}
	if gasInfo.MaxFee != 0 {
		txn.Fee.Amount = []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(gasInfo.MaxFee)},
		}
	} else if gasInfo.GasPrice != 0 {
		txn.Fee.Amount = []*TransactionFeeAmount{
			&TransactionFeeAmount{Denom: TOKEN_NAME, Amount: strconv.Itoa(feeGas * gasInfo.GasPrice)},
		}
	}

	// Set signatures
	if signature, err := transaction.Sign(txn); err != nil {
		return nil, err
	} else {
		txn.Signatures = []*TransactionSignature{
			&TransactionSignature{
				PubKey: &TransactionSignaturePubKey{
					Type:  tmsecp256k1.PubKeyAminoName,
					Value: base64.StdEncoding.EncodeToString(transaction.Client.privateKey.PubKey().SerializeCompressed()),
				},
				Signature:     signature,
				AccountNumber: strconv.Itoa(transaction.Client.account.AccountNumber),
				Sequence:      strconv.Itoa(transaction.Client.account.Sequence),
			},
		}
	}

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
	// transaction.Client.Infof("txn broadcast response %+v", string(body))
	// Read transaction broadcast response
	res := &TransactionBroadcastResponse{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	transaction.Client.Infof("txn broadcast response %+v", res)

	// https://github.com/bluzelle/blzjs/blob/45fe51f6364439fa88421987b833102cc9bcd7c0/src/swarmClient/cosmos.js#L240-L246
	// note - as of right now (3/6/20) the responses returned by the Cosmos REST interface now look like this:
	// success case: {"height":"0","txhash":"3F596D7E83D514A103792C930D9B4ED8DCF03B4C8FD93873AB22F0A707D88A9F","raw_log":"[]"}
	// failure case: {"height":"0","txhash":"DEE236DEF1F3D0A92CB7EE8E442D1CE457EE8DB8E665BAC1358E6E107D5316AA","code":4,
	//  "raw_log":"unauthorized: signature verification failed; verify correct account sequence and chain-id"}
	//
	// this is far from ideal, doesn't match their docs, and is probably going to change (again) in the future.

	if res.Code == 0 {
		transaction.Client.account.Sequence += 1
		if res.Data == "" {
			return []byte{}, nil
		}
		decodedData, err := hex.DecodeString(res.Data)
		return decodedData, err
	}
	if strings.Contains(res.RawLog, "signature verification failed") {
		transaction.broadcastRetries += 1
		transaction.Client.Warnf("transaction failed ... retrying(%d) ...", transaction.broadcastRetries)
		if transaction.broadcastRetries >= BROADCAST_MAX_RETRIES {
			return nil, fmt.Errorf("transaction failed after max retry attempts")
		}
		time.Sleep(BROADCAST_RETRY_INTERVAL)
		// Lookup changed sequence
		if err := transaction.Client.setAccount(); err != nil {
			return nil, err
		}
		b, err := transaction.Broadcast(txn)
		return b, err
	}

	return nil, fmt.Errorf("%s", res.RawLog)
}

func (transaction *Transaction) Sign(req *TransactionBroadcastPayload) (string, error) {
	payload := &TransactionBroadcastPayloadSignPayload{
		AccountNumber: strconv.Itoa(transaction.Client.account.AccountNumber),
		ChainId:       transaction.Client.options.ChainId,
		Sequence:      strconv.Itoa(transaction.Client.account.Sequence),
		Memo:          req.Memo,
		Fee:           req.Fee,
		Msgs:          req.Msg,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	transaction.Client.Infof("txn sign %+v", string(payloadBytes))
	hash := tmcrypto.Sha256(payloadBytes)
	if s, err := transaction.Client.privateKey.Sign(hash); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(serializeSig(s)), nil
	}
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
