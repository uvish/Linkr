package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	concurrency = 50
	requests    = 1000
)

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	jobs := make(chan int, requests)
	var writeSuccess, writeError, writeNonOK, readSuccess, readError, readNonOK uint64
	shortCodes := make(chan string, requests) // Channel to store short codes for read tests

	// Worker function
	worker := func() {
		defer wg.Done()
		for range jobs {
			// Write operation (POST /shorten)
			url := fmt.Sprintf("https://example.com/page/%d", rand.Intn(1000000))
			payload := map[string]string{"url": url}
			body, _ := json.Marshal(payload)

			start := time.Now()
			resp, err := http.Post("http://localhost:8080/shorten", "application/json", bytes.NewBuffer(body))
			if err != nil {
				atomic.AddUint64(&writeError, 1)
				fmt.Printf("Write request error: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				atomic.AddUint64(&writeNonOK, 1)
				fmt.Printf("Write non-OK response: %d\n", resp.StatusCode)
				continue
			}

			// Parse response to get short code
			var shortenResp ShortenResponse
			if err := json.NewDecoder(resp.Body).Decode(&shortenResp); err != nil {
				atomic.AddUint64(&writeError, 1)
				fmt.Printf("Write response parse error: %v\n", err)
				continue
			}

			atomic.AddUint64(&writeSuccess, 1)
			fmt.Printf("Write success: %vms\n", time.Since(start).Milliseconds())

			// Extract short code from short_url (e.g., "http://localhost:8080/abc123" -> "abc123")
			shortCode := shortenResp.ShortURL[strings.LastIndex(shortenResp.ShortURL, "/")+1:]
			shortCodes <- shortCode
		}
	}

	// Read worker function
	readWorker := func() {
		defer wg.Done()
		for shortCode := range shortCodes {
			start := time.Now()
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/%s", shortCode))
			if err != nil {
				atomic.AddUint64(&readError, 1)
				fmt.Printf("Read request error for %s: %v\n", shortCode, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusFound {
				atomic.AddUint64(&readNonOK, 1)
				fmt.Printf("Read non-OK response for %s: %d\n", shortCode, resp.StatusCode)
				continue
			}

			atomic.AddUint64(&readSuccess, 1)
			fmt.Printf("Read success for %s: %vms\n", shortCode, time.Since(start).Milliseconds())
		}
	}

	// Start workers
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	// Start read workers (use half the concurrency for reads)
	for i := 0; i < concurrency/2; i++ {
		wg.Add(1)
		go readWorker()
	}

	// Send write jobs
	for i := 0; i < requests; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for all write jobs to complete, then close shortCodes
	wg.Wait()
	close(shortCodes)

	// Wait for read workers to finish
	wg.Wait()

	fmt.Printf("Stress test completed in %vms\n", time.Since(start).Milliseconds())
	fmt.Printf("Write - Success: %d, Errors: %d, Non-OK: %d\n", writeSuccess, writeError, writeNonOK)
	fmt.Printf("Read  - Success: %d, Errors: %d, Non-OK: %d\n", readSuccess, readError, readNonOK)
}
