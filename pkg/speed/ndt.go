package speed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const ndtLocateURL = "https://locate.measurementlab.net/v2/nearest/ndt/ndt7"

// NDTLocateResult is the NDT locate API response structure
type NDTLocateResult struct {
	Results []NDTMachine `json:"results"`
}

// NDTMachine is an MeasurementLab machine hosting NDT
// The URLs are web socket addresses for connecting and
// testing against this machine
type NDTMachine struct {
	Machine  string      `json:"machine"`
	Location NDTLocation `json:"location"`
	URLs     []string    `json:"urls"`
}

// NDTLocation describes a rough location for an NDTMachine
type NDTLocation struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// Locate makes an NDT locate request, returning a list of
// the nearest NDTMachines that can be tested against
func Locate(client *http.Client) ([]NDTMachine, error) {
	var buf bytes.Reader
	req, err := http.NewRequest(http.MethodGet, ndtLocateURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("build ndt locate request: %v", err)
	}

	// Despite setting a do-not-track header, measurementlabs does
	// does still collect IP address and network diagnostic information
	// as per their privacy policy:
	// https://www.measurementlab.net/privacy/
	// I'm unsure how or if this changes measurementlab's behavior, but
	// chrome sends this header and it can't hurt to include.
	req.Header.Set("dnt", "1")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get ndt locate response: %v", err)
	}

	if res == nil || res.Body == nil {
		return nil, nil
	}

	var result NDTLocateResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode ndt locate response: %v", err)
	}

	return result.Results, nil
}
