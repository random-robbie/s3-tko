package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func main() {
	flag.Parse()

	var input io.Reader
	input = os.Stdin

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// Open log file once at the start
	logFile, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	sc := bufio.NewScanner(input)

	urls := make(chan string, 128)
	concurrency := 12
	var wg sync.WaitGroup
	var logMutex sync.Mutex
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for raw := range urls {
				func() {
					u, err := url.ParseRequestURI(raw)
					if err != nil {
						fmt.Fprintf(os.Stderr, "invalid url: %s\n", raw)
						return
					}

					if !resolves(u) {
						fmt.Fprintf(os.Stderr, "does not resolve: %s\n", u)
						return
					}

					resp, err := fetchURL(ctx, u)
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to fetch: %s (%s)\n", u, err)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						fmt.Printf("200 response code: %s (%s)\n", u, resp.Status)
					}
					if resp.StatusCode != http.StatusOK {
						fmt.Printf("response code: %s (%s)\n", u, resp.Status)
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							fmt.Fprintf(os.Stderr, "failed to read response body: %s\n", err)
							return
						}
						bodyStr := string(body)
						if strings.Contains(bodyStr, "<Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist</Message><BucketName>") {
							color.HiGreen("[*] Vulnerable System Found!\n")

							logMutex.Lock()
							if _, err := fmt.Fprintf(logFile, "%s/\n", u.String()); err != nil {
								log.Println(err)
							}
							logMutex.Unlock()
						}
					}
				}()
			}
		}()
	}

	for sc.Scan() {
		urls <- sc.Text()
	}
	close(urls)

	if sc.Err() != nil {
		fmt.Printf("error: %s\n", sc.Err())
	}

	wg.Wait()
}

func resolves(u *url.URL) bool {
	addrs, _ := net.LookupHost(u.Hostname())
	return len(addrs) != 0
}

func fetchURL(ctx context.Context, u *url.URL) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("User-Agent", "s3-tko scanner/0.1")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
