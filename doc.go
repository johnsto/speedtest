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

/*

Package speedtest is a client for the speedtest.net bandwidth measuring service.
This package is not affiliated, connected, or associated with speedtest.net in
any way. For information about speedtest.net, visit: http://www.speedtest.net/

This package is hosted on GitHub: http://www.github.com/johnsto/speedtest

This package implements only one part of speedtest.net's functionality, namely
the testing of your connection's upload and download bandwidth by way of HTTP
requests made to one of the service's many global testing servers. It does not
currently report the results back to the speedtest.net API.

Unlike speedtest.net, no attempt is made to ping servers to find the nearest
one. The 'nearest' server is selected using the geographic locations returned
by the speedtest.net API and may not always be physically correct.

A CLI is provided, which is the simplest and easiest way to measure your
connection's bandwidth.

Below is a simple example of how you might test bandwidth against the first
server reported by the service:

	import "fmt"
	import . "github.com/johnsto/speedtest"
	import "net/http"
	import "time"

	func main() {
		// Fetch server list
		settings, _ := FetchSettings()
		// Configure benchmark
		benchmark := NewDownloadBenchmark(http.DefaultClient, settings.Servers[0])
		// Run benchmark
		rate := RunBenchmark(benchmark, 4, 16, time.Second * 10)
		// Print result (bps)
		fmt.Println(NiceRate(rate))
	}

For a more detailed example, see speedtest-cli/cli.go

The algorithms used by this package differs from that used by the original
service. There are no guarantees about whether the approach used here is more
or less accurate, but my own measurements showed the results to be
broadly representative of line speed.

Both upload and download testing send/receive as much data as they can within
a given time period, increasing the number of concurrent requests as necessary.
Data is read/written in 16kb chunks, and the amount of data transferred is
logged into an array indexed by time, to a resolution of 100ms. Once all
transfers are complete, the array is scanned to find the second in
which cumulative data transfer was highest, and this value is combined with
the median transfer rate over the period to provide an estimated speed.

This approach is designed to produce a value close to the peak speed of the
line. In contrast, a mean average will typically underestimate due to the
overheads of the testing process.

*/
package speedtest
