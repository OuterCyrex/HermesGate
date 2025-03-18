package test

import (
	"net/http"
	"sync"
	"testing"
)

// BenchmarkConcurrentGetLocalhost 测试并发 GET 请求性能
func BenchmarkConcurrentGetLocalhost(b *testing.B) {
	client := &http.Client{}
	var wg sync.WaitGroup
	const maxConcurrency = 100
	sem := make(chan struct{}, maxConcurrency)

	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := client.Get("http://localhost:80/user")
			if err != nil {
				b.Errorf("Failed to send GET request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				var content []byte
				_, _ = resp.Body.Read(content)
				b.Errorf("Unexpected status code: %d, %s", resp.StatusCode, content)
			}
		}()
	}

	wg.Wait()
}
