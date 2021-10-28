package relay

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const apiUrl = "https://api.defender.openzeppelin.com"

type Client struct {
	poolId               string
	clientId             string
	region               string
	token                string
	username             string
	jsonRpcRequestNextId int64
}

func New(ctx context.Context, username string, password string) (*Client, error) {
	relay := &Client{
		poolId:               "us-west-2_iLmIggsiy",
		clientId:             "1bpd19lcr33qvg5cr3oi79rdap",
		region:               "us-west-2",
		username:             username,
		jsonRpcRequestNextId: 0,
	}

	token, err := relay.auth(ctx, username, password)
	if err != nil {
		return nil, err
	}
	relay.token = token

	return relay, nil
}

func (r *Client) ApiCall(ctx context.Context, method string, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, apiUrl+path, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Api-Key", r.username)
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyRes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: %s", resp.Status, string(bodyRes))
	}

	bodyRes, _ := ioutil.ReadAll(resp.Body)
	return bodyRes, nil
}

type JsonRpcRequest struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	JsonRpc string        `json:"jsonrpc"`
	Id      int64         `json:"id"`
}

func (r *Client) JsonRpc(ctx context.Context, method string, params []interface{}) ([]byte, error) {
	r.jsonRpcRequestNextId++
	body, err := json.Marshal(JsonRpcRequest{
		Method:  method,
		Params:  params,
		JsonRpc: "2.0",
		Id:      r.jsonRpcRequestNextId,
	})
	if err != nil {
		return nil, err
	}

	return r.ApiCall(ctx, "POST", "/relayer/jsonrpc", bytes.NewReader(body))
}

type Transaction struct {
	To         string `json:"to,omitempty"`
	GasLimit   string `json:"gasLimit,omitempty"`
	Data       string `json:"data,omitempty"`
	Speed      string `json:"speed,omitempty"`
	GasPrice   string `json:"gasPrice,omitempty"`
	Value      string `json:"value,omitempty"`
	ValidUntil string `json:"validUntil,omitempty"`
}

func (r *Client) SendTransaction(ctx context.Context, tx *Transaction) ([]byte, error) {
	body, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return r.ApiCall(ctx, "POST", "/txs", bytes.NewReader(body))
}

func (r *Client) ReplaceTransactionById() {
	// TODO
}

func (r *Client) ReplaceTransactionByNonce() {
	// TODO
}

func (r *Client) Sign(ctx context.Context, message string) ([]byte, error) {
	body, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "0x" + hex.EncodeToString([]byte(message)),
	})
	if err != nil {
		return nil, err
	}

	return r.ApiCall(ctx, "POST", "/sign", bytes.NewReader(body))
}

func (r *Client) Query(ctx context.Context, id string) ([]byte, error) {
	return r.ApiCall(ctx, "GET", "/query/"+id, nil)
}

type TransactionInfo struct {
	ChainID       int       `json:"chainId"`
	Hash          string    `json:"hash"`
	TransactionID string    `json:"transactionId"`
	Value         string    `json:"value"`
	GasPrice      int64     `json:"gasPrice"`
	GasLimit      int       `json:"gasLimit"`
	To            string    `json:"to"`
	From          string    `json:"from"`
	Data          string    `json:"data"`
	Nonce         int       `json:"nonce"`
	Status        string    `json:"status"`
	Speed         string    `json:"speed"`
	ValidUntil    time.Time `json:"validUntil"`
	CreatedAt     time.Time `json:"createdAt"`
	SentAt        time.Time `json:"sentAt"`
	PricedAt      time.Time `json:"pricedAt"`
}

func (r *Client) List(ctx context.Context) ([]TransactionInfo, error) {
	// TODO criteria
	var txs []TransactionInfo

	result, err := r.ApiCall(ctx, "GET", "/txs", nil)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(result, &txs); err != nil {
		return nil, err
	}

	return txs, nil
}
