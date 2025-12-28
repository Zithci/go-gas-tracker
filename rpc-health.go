package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RPCHealth struct {
	Name     string
	URL      string
	Status   string
	Duration time.Duration
}

func checkRPC(ctx context.Context, name, url string, wg *sync.WaitGroup, results chan<- RPCHealth) {
	defer wg.Done()

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		results <- RPCHealth{Name: name, URL: url, Status: "ERROR", Duration: 0}
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	elapsed := time.Since(start)

	if err != nil {
		if err == context.DeadlineExceeded {
			results <- RPCHealth{Name: name, URL: url, Status: "TIMEOUT", Duration: elapsed}
		} else {
			results <- RPCHealth{Name: name, URL: url, Status: "ERROR", Duration: elapsed}
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results <- RPCHealth{Name: name, URL: url, Status: "HEALTHY", Duration: elapsed}
	} else {
 status := fmt.Sprintf("UNHEALTHY (%d)", resp.StatusCode)
    results <- RPCHealth{Name: name, URL: url, Status: status, Duration: elapsed}
	


}
}

func main() {
	rpcs := map[string]string{
		"Ethereum":  "https://eth.llamarpc.com",
		"Arbitrum":  "https://arb1.arbitrum.io/rpc",
		"Base":      "https://mainnet.base.org",
		"Optimism":  "https://mainnet.optimism.io",
		"Polygon":   "https://polygon-rpc.com",
		"FakeRPC":   "https://this-will-timeout.fake",
	}

	var wg sync.WaitGroup
	results := make(chan RPCHealth, len(rpcs))

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  RPC HEALTH CHECK (3s timeout)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for name, url := range rpcs {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		wg.Add(1)
		go checkRPC(ctx, name, url, &wg, results)
	}

	wg.Wait()
	close(results)

	healthy := 0
	unhealthy := 0
	timeout := 0

	for result := range results {
		status := result.Status
		icon := "✅"
		if status == "TIMEOUT" {
			icon = "⏰"
			timeout++
		} else if status == "ERROR" || status == "UNHEALTHY" {
			icon = "❌"
			unhealthy++
		} else {
			healthy++
		}

		fmt.Printf("%-12s %s %-10s (%v)\n", result.Name, icon, status, result.Duration.Round(time.Millisecond))
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Healthy: %d | Unhealthy: %d | Timeout: %d\n", healthy, unhealthy, timeout)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}