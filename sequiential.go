package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	// "start"
)

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

	type RPCResponse struct {
		JSONRPC string `json:"jsonrpc"`
	ID int 	`json:"id"`
	Result string `json:"result"`
}

func getGasPrice (rpcURL string)(float64, error) {
	// RPC req

	reqBody := RPCRequest{
		JSONRPC : "2.0",
		Method : "eth_gasPrice",
		Params: []interface{}{},
		ID: 1,
	}

	jsonData, _ := json.Marshal(reqBody)

	//send POST Request
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close() //close http req

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	//parse response
	var rpcResp RPCResponse
	json.Unmarshal(body, &rpcResp)

	//hex -> decimal
	hexStr := strings.TrimPrefix(rpcResp.Result, "0x")
	weiPrice, _ :=strconv.ParseInt(hexStr, 16,64)


	//wei -> gwei (divide by 1e9)
	gweiPrice := float64(weiPrice) /1e9

	return gweiPrice,nil
}

func main() {
	start := time.Now()
	//header
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  GAS PRICES (Gwei-sequiental)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")

	//ETH
	ethRPC :=  "https://eth.llamarpc.com"
	gasPrice, err := getGasPrice(ethRPC)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	fmt.Printf("Ethereum gas: %.2f gwei\n", gasPrice)
	
	//ARBITRUM
	arbRPC :=  "https://arb1.arbitrum.io/rpc"
	arbGas,err := getGasPrice(arbRPC)
	if err != nil {
		fmt.Println("Arbitrum err:", err)
	}else {
		fmt.Printf("Arbitrum gas: %.2f gwei\n", arbGas)
	}

	//BASE
	baseRPC := "https://mainnet.base.org"
	baseGas, err := getGasPrice(baseRPC)
	if err != nil {
		fmt.Println("Base Error:", err)
	} else {
		fmt.Printf("Base gas:%.2f gwei\n", baseGas)
	}

	//footer v
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")


	//run time 
	elapsed := time.Since(start)  // ← Add this
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("⏱️  Time: %v\n", elapsed)  // ← Add this
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")
}