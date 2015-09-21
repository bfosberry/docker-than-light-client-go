package dtl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

var ErrScanFailed = errors.New("Scan failed!")
var ErrFireFailed = errors.New("Fire failed!")
var ErrTravelFailed = errors.New("Travel failed!")

type Client interface {
	ScanSector() ([]*Ship, []*Sector, *Ship, error)
	Travel(*Sector) (*Ship, error)
	Fire(string) (*Ship, error)
}

type client struct {
	ApiURL string
	token  string
}

type ScanSectorResponse struct {
	Ships   []*Ship   `json:"ships"`
	Sectors []*Sector `json:"sectors"`
	State   *Ship     `json:"state"`
}

type TravelResponse struct {
	State *Ship
}

type FireResponse struct {
	State *Ship
}

func (c *client) ScanSector() ([]*Ship, []*Sector, *Ship, error) {
	var r ScanSectorResponse
	if err := c.makeRequest("GET", "scan", &r, ErrScanFailed); err != nil {
		return nil, nil, nil, err
	}
	return r.Ships, r.Sectors, r.State, nil
}

func (c *client) Travel(sector *Sector) (*Ship, error) {
	var r TravelResponse
	if err := c.makeRequest("POST", fmt.Sprintf("travel/%s", sector.Name), &r, ErrTravelFailed); err != nil {
		return nil, err
	}

	return r.State, nil
}

func (c *client) Fire(target string) (*Ship, error) {
	var r FireResponse
	if err := c.makeRequest("POST", fmt.Sprintf("fire/%s", target), &r, ErrFireFailed); err != nil {
		return nil, err
	}
	return r.State, nil
}

func (c *client) makeRequest(method, url string, r interface{}, defaultError error) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.ApiURL, url), nil)
	if err != nil {
		return err
	}
	req.Header.Set("ContentType", "application/json")
	req.Header.Set("Authorization", c.token)

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return defaultError
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		return err
	}
	return nil
}

func NewClient() (Client, error) {
	apiURL := os.Getenv("API_URL")
	token := os.Getenv("TOKEN")
	if apiURL == "" {
		return nil, errors.New("No API url provided")
	}
	if token == "" {
		return nil, errors.New("No Token provided")
	}
	return &client{
		ApiURL: apiURL,
		token:  token,
	}, nil
}
