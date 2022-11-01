package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	const lookFor = "Go"

	var (
		workersLimit = 5
		urlsList     = []string{
			"https://golang.org",
			"https://market.yandex.ru",
			"https://www.youtube.com",
			"https://github.com",
			"https://golang.org",
			"https://golang.org",
			"https://golang.org",
			"https://golang.org",
		}

		total = 0
	)

	if len(urlsList) < workersLimit {
		workersLimit = len(urlsList)
	}

	ctx, cancel := context.WithCancel(context.Background())

	urlCh := make(chan string)
	go func() {
		defer close(urlCh)
		for _, url := range urlsList {
			urlCh <- url
		}
	}()
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < workersLimit; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			for r := range urlCh {
				body, err := makeRequest(r)
				if err != nil {
					log.Fatalln(err)
				}
				contain := strings.Count(string(body), lookFor)
				fmt.Printf("Count for %s: %d\n", r, contain)
				mu.Lock()
				total += contain
				mu.Unlock()

			}
		}(ctx)
	}

	wg.Wait()
	fmt.Printf("Total: %d\n", total)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	<-sigterm
	// graceful shutdown
	log.Print("graceful shutdown...")
	cancel()
}

func makeRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	return io.ReadAll(resp.Body)
}
