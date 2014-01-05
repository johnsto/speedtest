/*
The MIT License (MIT)

Copyright (c) 2014 David Johnston

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package speedtest

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	chunkSize  = 16384
	windowSize = 10
)

// ErrTimeExpired is returned by readers/writers if they were halted due to
// exceeding the duration of the benchmark.
var ErrTimeExpired = errors.New("time expired")

// A Benchmark represents a specific bandwidth test.
type Benchmark interface {
	Run(func(n int) error) error
}

// DownloadBenchmark represents a download bandwidth test.
type DownloadBenchmark struct {
	Client  http.Client
	Server  Server
	BaseURL string
}

// NewDownloadBenchmark creates a new download benchmark with the given HTTP
// client and test server.
func NewDownloadBenchmark(client http.Client, server Server) DownloadBenchmark {
	slashPos := strings.LastIndex(server.URL, "/")
	baseURL := server.URL[:slashPos] + "/random1000x1000.jpg"
	return DownloadBenchmark{client, server, baseURL}
}

// Run fetches a file, reporting the size of each downloaded chunk to the
// callback function, ending only on EOF or when the callback returns an error.
func (b DownloadBenchmark) Run(fn func(n int) error) error {
	threadURL := b.BaseURL + "?x=" + strconv.Itoa(rand.Int())
	resp, err := b.Client.Get(threadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := make([]byte, chunkSize)
	for {
		num, err := resp.Body.Read(buf)
		nerr := fn(num)
		if nerr == ErrTimeExpired || err == io.EOF {
			break
		}
		if nerr != nil {
			return nerr
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadBenchmark represents an upload bandwidth test.
type UploadBenchmark struct {
	Client http.Client
	Server Server
}

// NewUploadBenchmark creates a new upload benchmark with the given HTTP
// client and test server.
func NewUploadBenchmark(client http.Client, server Server) UploadBenchmark {
	return UploadBenchmark{client, server}
}

// Run performs an HTTP POST, uploading junk data and reporting the size of
// each uploaded chunk.
func (b UploadBenchmark) Run(fn func(n int) error) error {
	reader := NewJunkReader(1024 * 1024)
	writer := NewCallbackWriter(fn)
	tee := io.TeeReader(&reader, writer)
	_, err := b.Client.Post(b.Server.URL, "text/plain", tee)
	return err
}

// RunBenchmark runs the given benchmark for the given amount of time. It
// increases the number of threads up to the maximum as each one finishes.
// The returned value is the maximum number of bytes recorded from any
// contiguous 1 second window within the testing period.
func RunBenchmark(b Benchmark, threads int, maxThreads int, duration time.Duration) int {
	var wg sync.WaitGroup
	var tc sync.Mutex

	// Setup sampling parameters
	resolution := time.Second / time.Duration(windowSize)
	chunks := make([]int, duration/resolution)

	// Seed thread pool
	reqs := make(chan int, maxThreads)
	for i := 0; i < threads; i++ {
		reqs <- 1
	}

	// Setup timeout
	start := time.Now()
	active := true

	perform := func() {
		wg.Add(1)
		defer wg.Done()

		// Run benchmark, recording reads into timestamped array
		err := b.Run(func(n int) error {
			p := int(time.Since(start) / resolution)
			if p < len(chunks) {
				chunks[p] += n
			}
			if !active {
				return ErrTimeExpired
			}
			return nil
		})

		if active {
			if err != nil {
				log.Fatalln(err)
			}

			// Enqueue next task
			reqs <- 1

			// See if we can add another thread
			tc.Lock()
			if threads < maxThreads {
				threads++
				reqs <- 1
			}
			tc.Unlock()
		} else {
			// All done!
			return
		}
	}

	// Process queue
	timeout := time.After(duration)
	for active {
		select {
		case <-reqs:
			if active {
				go perform()
			}
		case <-timeout:
			// Outta time, signal
			active = false
		}
	}

	wg.Wait()

	maxSum := MaximalSumWindow(chunks, windowSize)
	windowAvg := MedianSumWindow(chunks, windowSize)
	return (maxSum + windowAvg) / 2
}
