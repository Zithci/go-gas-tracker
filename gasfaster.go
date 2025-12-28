package main
 
import(
   	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)


// slot buat chain
type ChainStats struct {
	Name string
	Min float64
	Max float64
	Sum float64
	Count int
}

// create slot buat placement data dari
type RPCRequest struct {
	JSONRPC string  `json:"jsonrpc"`
	Method string `json:"method"`
	Params []interface{} `json:"params"`
	ID int `json:"id"`
}

type RPCResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	ID int `json:"id"`
	Result string `json:"result"`
}

// handle slot gas
type ChainGas struct{
	 Name string 
	 Price float64
	 Error error 
}

func getGasPrice(rpcURL string) (float64, error){
	reqBody := RPCRequest{
		JSONRPC : "2.0",
		Method : "eth_gasPrice",
		Params : []interface{}{},
		ID: 1,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
	return 0, err
	}
	defer resp.Body.Close()

	//err handle
  body, err := io.ReadAll(resp.Body)
	if err != nil {
	return 0, err
	}
	
	var rpcResp RPCResponse
	json.Unmarshal(body, &rpcResp)

	hexStr := strings.TrimPrefix(rpcResp.Result, "0x")
	weiPrice, _ := strconv.ParseInt(hexStr, 16, 64)
	gweiPrice := float64(weiPrice) / 1e9

	return gweiPrice, nil
}

func fetchGasParallel(name string, rpcURL string, wg *sync.WaitGroup, results chan<- ChainGas) {
	defer wg.Done()

	price, err := getGasPrice(rpcURL)
	results <- ChainGas{
		Name:  name,
		Price: price,
		Error: err,
	}
}

func main() {
	start := time.Now()
	chains := map[string]string{
		"Ethereum": "https://eth.llamarpc.com",
		"Arbitrum": "https://arb1.arbitrum.io/rpc",
		"Base":     "https://mainnet.base.org",
		"Optimism": "https://mainnet.optimism.io",
		"Polygon":  "https://polygon-rpc.com",
	}

// Stats tracking
	stats := make(map[string]*ChainStats)
	for name := range chains {
		stats[name] = &ChainStats{
			Name:  name,
			Min:   999999,
			Max:   0,
			Sum:   0,
			Count: 0,
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("       LIVE GAS PRICE TRACKER")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Loop 5 updates
	for i := 1; i <= 3; i++ {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("📊 Update #%d [%s]\n", i, timestamp)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		var wg sync.WaitGroup
		results := make(chan ChainGas, len(chains))

		for name, rpc := range chains {
			wg.Add(1)
			go fetchGasParallel(name, rpc, &wg, results)
		}

		wg.Wait()
		close(results)

		for result := range results {
			if result.Error != nil {
				fmt.Printf("%-12s ❌ ERROR\n", result.Name)
			} else {
				fmt.Printf("%-12s %.2f gwei\n", result.Name, result.Price)

				// Update stats
				s := stats[result.Name]
				s.Count++
				s.Sum += result.Price
				if result.Price < s.Min {
					s.Min = result.Price
				}
				if result.Price > s.Max {
					s.Max = result.Price
				}
			}
		}

		fmt.Println()

		if i < 5 {
			time.Sleep(3 * time.Second)
		}
	}

	// Print stats
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("       STATISTICS (3 updates)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, s := range stats {
		if s.Count > 0 {
			avg := s.Sum / float64(s.Count)
			fmt.Printf("%-12s Min: %6.2f | Max: %6.2f | Avg: %6.2f\n",
				s.Name, s.Min, s.Max, avg)
		}
	}
	elapsed := time.Since(start)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("⏱️  Total Time:  %.1fs\n", elapsed.Seconds())
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}





