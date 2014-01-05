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

package main

import (
	"flag"
	"fmt"
	"github.com/johnsto/speedtest"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	findNearest = -1
	findFarthest = -2
)

var (
	cmdListServers   string
	testUpload       bool
	testDownload     bool
	httpTimeout      time.Duration
	sampleServer     int
	samplePeriod     time.Duration
	sampleThreads    int
	sampleMaxThreads int
)

func init() {
	flag.StringVar(&cmdListServers, "list-servers", "",
		"List servers (id|distance|nearest|farthest)")

	flag.BoolVar(&testUpload, "test-upload", true, "Test upload speed")
	flag.BoolVar(&testDownload, "test-download", true, "Test download speed")

	flag.IntVar(&sampleServer, "server", findNearest,
		"Server id to test (-1: use nearest, -2: use farthest)")

	flag.DurationVar(&httpTimeout, "timeout", time.Duration(10*time.Second),
		"HTTP connection timeout")
	flag.DurationVar(&samplePeriod, "period", time.Duration(10*time.Second),
		"Sampling period")
	flag.IntVar(&sampleThreads, "threads", 4,
		"Initial number of benchmark threads")
	flag.IntVar(&sampleMaxThreads, "max-threads", 16,
		"Maximum number of benchmark threads")
}

func main() {
	flag.Parse()

	fmt.Printf("Fetching server list... ")
	settings, err := speedtest.FetchSettings()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%v found.\n", len(settings.Servers))

	fmt.Printf("Fetching config...\n")
	config, err := speedtest.FetchConfig()
	if err != nil {
		fmt.Printf("Couldn't read config: %v", err)
		os.Exit(1)
	}
	settings.UpdateDistances(config.Client.Lat, config.Client.Lon)

	fmt.Printf("  ISP: %v\n", config.Client.IspName)
	fmt.Printf("  Location: %v, %v\n\n", config.Client.Lat, config.Client.Lon)

	// List servers
	if cmdListServers != "" {
		var listing = settings.Servers
		switch cmdListServers {
		case "id":
			settings.Servers.SortByID()
		case "distance":
			settings.Servers.SortByDistance()
		case "nearest":
			settings.Servers.SortByDistance()
			listing = settings.Servers[:10]
		case "farthest":
			settings.Servers.SortByDistance()
			listing = settings.Servers[len(listing)-10:]
		}
		for _, server := range listing {
			fmt.Printf("%5d. [%v] (%dkm) %v\n",
				server.ID, server.CountryCode, int(server.Distance), server.Name)
		}
		return
	}

	var server speedtest.Server
	switch sampleServer {
	case findNearest:
		settings.Servers.SortByDistance()
		server = settings.Servers[0]
	case findFarthest:
		settings.Servers.SortByDistance()
		server = settings.Servers[len(settings.Servers)-1]
	default:
		// find server with ID
		for _, s := range settings.Servers {
			if s.ID == sampleServer {
				server = s
				break
			}
		}
	}

	if server.ID == 0 {
		fmt.Printf("Could not find server. Re-run with -list-servers for a list.\n")
		os.Exit(1)
	}

	fmt.Printf("Using server %d. %v, %v, %v (%dkm)\n",
		server.ID, server.Sponsor, server.Name, server.Country, int(server.Distance))

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, httpTimeout)
			},
		},
	}

	if testDownload {
		benchmark := speedtest.NewDownloadBenchmark(client, server)
		fmt.Print("Testing download speed... ")
		rate := speedtest.RunBenchmark(benchmark, sampleThreads, sampleMaxThreads, samplePeriod)
		fmt.Println(speedtest.NiceRate(rate))
	}

	if testUpload {
		benchmark := speedtest.NewUploadBenchmark(client, server)
		fmt.Printf("Testing upload speed... ")
		rate := speedtest.RunBenchmark(benchmark, sampleThreads, sampleMaxThreads, samplePeriod)
		fmt.Println(speedtest.NiceRate(rate))
	}
}
