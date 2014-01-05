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
	"sort"
)

// Server encompasses a server definition returned by the speedtest.net API
type Server struct {
	ID          int     `xml:"id,attr"`
	URL         string  `xml:"url,attr"`
	Lat         float64 `xml:"lat,attr"`
	Lon         float64 `xml:"lon,attr"`
	Name        string  `xml:"name,attr"`
	Country     string  `xml:"country,attr"`
	CountryCode string  `xml:"cc,attr"`
	Sponsor     string  `xml:"sponsor,attr"`
	// Distance is calculated locally from the client configuration
	Distance float64
}

// Servers represents a sortable list of servers
type Servers []Server

// Len returns the length
func (s Servers) Len() int      { return len(s) }
// Swap swaps the items at the specified indexes
func (s Servers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// byDistance sorts servers by their measured Distance
type byDistance struct{ Servers }
func (s byDistance) Less(i, j int) bool {
	return s.Servers[i].Distance < s.Servers[j].Distance
}

// byID sorts servers by their ID
type byID struct{ Servers }
func (s byID) Less(i, j int) bool {
	return s.Servers[i].ID < s.Servers[j].ID
}

// SortByID sorts their servers by their numerical ID.
func (s Servers) SortByID() {
	sort.Sort(byID{s})
}

// SortByDistance sorts the servers by their measured distance.
func (s Servers) SortByDistance() {
	sort.Sort(byDistance{s})
}

// Settings encompasses server settings data provided by the speedtest.net API
type Settings struct {
	XMLName xml.Name `xml:"settings"`
	Servers Servers  `xml:"servers>server"`
}

// Config encompasses configuration data provided by the speedtest.net API
type Config struct {
	XMLName xml.Name `xml:"settings"`
	Client  Client   `xml:"client"`
}

// Client encompasses client information provided by the speedtest.net API
type Client struct {
	IPAddress string  `xml:"ip,attr"`
	Lat       float64 `xml:"lat,attr"`
	Lon       float64 `xml:"lon,attr"`
	IspName   string  `xml:"isp,attr"`
}

// UpdateDistances updates the Servers with the current latitude/longitude
func (settings Settings) UpdateDistances(lat float64, lon float64) {
	for i, server := range settings.Servers {
		settings.Servers[i].Distance = Distance(
			server.Lat*degToRad, server.Lon*degToRad,
			lat*degToRad, lon*degToRad)
	}
}
