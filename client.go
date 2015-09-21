package dtl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Client interface {
	ScanSector() ([]*Ship, []*Sector, *Ship, error)
	Travel(*Sector) (*Ship, error)
	Fire(string) (*Ship, error)
}

type client struct {
	ApiURL string
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
	resp, err := http.Get(fmt.Sprintf("%s/scan", c.ApiURL))
	if err != nil {
		return nil, nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, nil, errors.New("Scan failed!")
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		return nil, nil, nil, err
	}
	return r.Ships, r.Sectors, r.State, nil
}

func (c *client) Travel(sector *Sector) (*Ship, error) {
	var r TravelResponse
	resp, err := http.Post(fmt.Sprintf("%s/travel/%s", c.ApiURL, sector.Name), "application/json", nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Travel failed!")
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.State, nil
}

func (c *client) Fire(target string) (*Ship, error) {
	var r FireResponse
	resp, err := http.Post(fmt.Sprintf("%s/fire/%s", c.ApiURL, target), "application/json", nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Firing failed!")
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.State, nil
}

func NewClient() (Client, error) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		return nil, errors.New("No API url provided")
	}
	return &client{
		ApiURL: apiURL,
	}, nil
}
