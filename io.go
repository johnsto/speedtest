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

import "io"
import "fmt"

// A JunkReader produces junk-ish data
type JunkReader struct {
	Data []byte
	Size int
	Pos  int
}

// NewJunkReader creates a JunkReader that can be used to generate junk data
// of the specified size. A negative size is unbounded.
func NewJunkReader(size int) JunkReader {
	return JunkReader{
		Size: size,
		Pos:  0,
	}
}

// Read reads junk data into the specified buffer until the input buffer is
// filled or the reader is exhausted of data.
func (r *JunkReader) Read(p []byte) (n int, err error) {
	for {
		if r.Size >= 0 && r.Pos >= r.Size {
			// Outta data
			return n, io.EOF
		} else if n < len(p) {
			p[n] = byte(n)
			n++
			r.Pos++
		} else {
			// Outta buffer space
			return n, err
		}
	}
}

// A CallbackWriter is a Writer that calls a given callback for each Write. If
// the callback func returns an error, this is bubbled up from Write.
type CallbackWriter struct {
	Callback func(n int) error
}

// NewCallbackWriter creates a CallbackWriter with the specified callback function.
func NewCallbackWriter(callback func(n int) error) CallbackWriter {
	return CallbackWriter{callback}
}

// Write writes data to the writer, returning the length of the output buffer,
// and any error returned by the callback.
func (b CallbackWriter) Write(p []byte) (n int, err error) {
	return len(p), b.Callback(len(p))
}

// NiceRate represents a value measured in bytes/sec as bps, kbps or mbps as
// appropriate.
func NiceRate(rate int) string {
	bps := float64(8 * rate)
	kbps := bps / 1024
	mbps := kbps / 1024

	if mbps > 0.1 {
		return fmt.Sprintf("%.2fmbps", mbps)
	} else if kbps > 0.1 {
		return fmt.Sprintf("%.2fkbps", kbps)
	} else {
		return fmt.Sprintf("%.2fbps", bps)
	}
}
