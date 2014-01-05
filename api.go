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
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Fetch GETs a URL and returns the response body
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// FetchSettings fetches the list of available servers
func FetchSettings() (Settings, error) {
	body, err := Fetch("http://www.speedtest.net/speedtest-servers.php")
	if err != nil {
		return Settings{}, err
	}
	settings := Settings{}
	err = xml.Unmarshal(body, &settings)
	return settings, err
}

// FetchConfig fetches the recommended client configuration
func FetchConfig() (Config, error) {
	body, err := Fetch("http://www.speedtest.net/speedtest-config.php")
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = xml.Unmarshal(body, &config)
	return config, err
}
