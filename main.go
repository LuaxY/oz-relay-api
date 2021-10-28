package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"oz-relay-api/relay"
)

func main() {
	// --key=RELAY_API_KEY --secret=RELAY_API_SECRET
	apiKey := flag.String("key", "", "OZ Relay API key")
	apiSecret := flag.String("secret", "", "OZ Relay API secret")
	flag.Parse()

	ctx := context.Background()

	client, err := relay.New(ctx, *apiKey, *apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	//ExampleSendTx(ctx, client)
	//ExampleJsonRpc(ctx, client)
	ExampleSignMessage(ctx, client)
	//ExampleQueryTx(ctx, client)
	//ExampleListTxs(ctx, client)
}

func ExampleSendTx(ctx context.Context, client *relay.Client) {
	result, err := client.SendTransaction(ctx, &relay.Transaction{
		To:       "0x145d67831F4ce5F1989a02aA8aF6769f95d2A8d4",
		GasLimit: "0x100590",
		Data:     "0x47153f82000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000007f686f9b2740e09469a4a5a0f2ede2d317d4aee50000000000000000000000001a18f04e0aed1f1ab9f5cf1955d4f04d5c49f3e2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f4240000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000006423b872dd0000000000000000000000007f686f9b2740e09469a4a5a0f2ede2d317d4aee5000000000000000000000000dead00000000000000004206942069420694206900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000414b5485b0dcf7552761ccb37dc268e21ac31671ffb4fb86196200adf9da47d2d170f9a61463ac0263e719d91c8347e8e762e46c8ab1b06fe4ecbf51789e87bf231b00000000000000000000000000000000000000000000000000000000000000",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}

func ExampleJsonRpc(ctx context.Context, client *relay.Client) {
	var params []interface{}
	params = append(params, struct {
		From string `json:"from"`
		To   string `json:"to"`
		Data string `json:"data"`
		Gas  string `json:"gas"`
	}{
		From: "0x99e02ccd4146714b8bd7a1f818fc88e892658650",
		To:   "0x145d67831f4ce5f1989a02aa8af6769f95d2a8d4",
		Data: "0xbf5d3bdb000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000007f686f9b2740e09469a4a5a0f2ede2d317d4aee50000000000000000000000001a18f04e0aed1f1ab9f5cf1955d4f04d5c49f3e2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f4240000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000006423b872dd0000000000000000000000007f686f9b2740e09469a4a5a0f2ede2d317d4aee5000000000000000000000000dead00000000000000004206942069420694206900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000414b5485b0dcf7552761ccb37dc268e21ac31671ffb4fb86196200adf9da47d2d170f9a61463ac0263e719d91c8347e8e762e46c8ab1b06fe4ecbf51789e87bf231b00000000000000000000000000000000000000000000000000000000000000",
		Gas:  "0xaffaf8",
	})
	params = append(params, "latest")

	result, err := client.JsonRpc(ctx, "eth_call", params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}

func ExampleSignMessage(ctx context.Context, client *relay.Client) {
	result, err := client.Sign(ctx, "hello")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}

func ExampleQueryTx(ctx context.Context, client *relay.Client) {
	result, err := client.Query(ctx, "6ae57998-f56d-42e6-80a1-acae88706c44")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}

func ExampleListTxs(ctx context.Context, client *relay.Client) {
	transactions, err := client.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", transactions)
}