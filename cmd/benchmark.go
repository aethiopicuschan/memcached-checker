package cmd

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run benchmarks for Set, Get, and Del.",
	Long:  "Run benchmarks for Set, Get, and Del.",
	RunE:  benchmark,
}

func init() {
	benchmarkCmd.Flags().StringP("address", "a", "127.0.0.1:11211", "Address of the memcached server")
	benchmarkCmd.Flags().IntP("concurrency", "c", 10, "Number of concurrent workers")
	benchmarkCmd.Flags().IntP("requests", "r", 10000, "Request count per worker")

	rootCmd.AddCommand(benchmarkCmd)
}

type opLatency struct {
	set time.Duration
	get time.Duration
	del time.Duration
}

func benchmarkWorker(wg *sync.WaitGroup, address string, requests int, resultChan chan<- opLatency) {
	defer wg.Done()
	mc := memcache.New(address)
	defer mc.Close()

	for i := range requests {
		key := fmt.Sprintf("key_%s", uuid.New())
		value := fmt.Sprintf("value_%d", i)
		startSet := time.Now()
		if err := mc.Set(&memcache.Item{Key: key, Value: []byte(value)}); err != nil {
			fmt.Printf("Set error: %v\n", err)
			continue
		}
		setLatency := time.Since(startSet)

		startGet := time.Now()
		v, err := mc.Get(key)
		if err != nil {
			fmt.Printf("Get error: %v\n", err)
			continue
		}
		if string(v.Value) != value {
			fmt.Printf("Get value mismatch: expected %s, got %s\n", value, string(v.Value))
			continue
		}
		getLatency := time.Since(startGet)

		startDel := time.Now()
		if err := mc.Delete(key); err != nil {
			fmt.Printf("Delete error: %v\n", err)
			continue
		}
		delLatency := time.Since(startDel)

		resultChan <- opLatency{set: setLatency, get: getLatency, del: delLatency}
	}
}

func benchmark(cmd *cobra.Command, args []string) (err error) {
	concurrency, err := strconv.Atoi(cmd.Flag("concurrency").Value.String())
	if err != nil {
		return fmt.Errorf("invalid concurrency value: %v", err)
	}
	requests, err := strconv.Atoi(cmd.Flag("requests").Value.String())
	if err != nil {
		return fmt.Errorf("invalid requests value: %v", err)
	}
	totalRequests := concurrency * requests
	address := cmd.Flag("address").Value.String()
	fmt.Printf("Start benchmark: target=%s, workers=%d, request count per worker=%d (total request count=%d)\n", address, concurrency, requests, totalRequests)

	var wg sync.WaitGroup
	resultChan := make(chan opLatency, totalRequests)
	startTime := time.Now()

	for range concurrency {
		wg.Add(1)
		go benchmarkWorker(&wg, address, requests, resultChan)
	}
	wg.Wait()
	close(resultChan)
	totalTime := time.Since(startTime)

	var totalSetLatency time.Duration
	var totalGetLatency time.Duration
	var totalDelLatency time.Duration
	var count int
	for op := range resultChan {
		totalSetLatency += op.set
		totalGetLatency += op.get
		totalDelLatency += op.del
		count++
	}

	avgSetLatency := totalSetLatency / time.Duration(count)
	avgGetLatency := totalGetLatency / time.Duration(count)
	avgDelLatency := totalDelLatency / time.Duration(count)

	setQPS := float64(count) / totalTime.Seconds()
	getQPS := float64(count) / totalTime.Seconds()
	delQPS := float64(count) / totalTime.Seconds()

	fmt.Println("------ Benchmark Result ------")
	fmt.Printf("Total time: %v\n", totalTime)
	fmt.Printf("Total request: %d\n", count)
	fmt.Printf("Set - average latency: %v, QPS: %.2f\n", avgSetLatency, setQPS)
	fmt.Printf("Get - average latency: %v, QPS: %.2f\n", avgGetLatency, getQPS)
	fmt.Printf("Del - average latency: %v, QPS: %.2f\n", avgDelLatency, delQPS)

	return
}
