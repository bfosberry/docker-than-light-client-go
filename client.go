package dtl

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmcvetta/napping"
)

type Client interface {
	ScanSector() ([]Ship, []Sector, Ship, error)
	Travel(Sector) (Ship, error)
	Fire(string) (Ship, error)
}

type client struct {
	ApiURL string
}

type ScanSectorResponse struct {
	Ships   []Ship
	Sectors []Sector
	NewShip Ship
}

type TravelResponse struct {
	NewShip Ship
}

type FireResponse struct {
	NewShip Ship
}

func (c *client) ScanSector() ([]Ship, []Sector, Ship, error) {
	s := napping.Session{}
	var errString string
	r := ScanSectorResponse{}
	resp, err := s.Get(fmt.Sprintf("%s/scan", c.ApiURL), nil, &r, &errString)

	if err != nil {
		return nil, nil, nil, err
	}
	if resp.Status() != 200 {
		return nil, nil, nil, errors.New("Scan failed!")
	}
	return r.Ships, r.Sectors, r.NewShip, nil
}

func (c *client) Travel(sector Sector) (Ship, error) {
	s := napping.Session{}
	var errString string
	r := TravelResponse{}
	resp, err := s.Get(fmt.Sprintf("%s/travel/%s", c.ApiURL, sector.Name), nil, &r, &errString)

	if err != nil {
		return nil, err
	}
	if resp.Status() != 200 {
		return nil, errors.New("Scan failed!")
	}
	return r.NewShip, nil
}

func (c *client) Fire(target string) (Ship, error) {
	s := napping.Session{}
	var errString string
	r := FireResponse{}
	resp, err := s.Get(fmt.Sprintf("%s/fire/%s", c.ApiURL, target), nil, &r, &errString)

	if err != nil {
		return nil, err
	}
	if resp.Status() != 200 {
		return nil, errors.New("Scan failed!")
	}
	return r.NewShip, nil
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
