package dtl

import (
	"encoding/json"
	"io"
	"os"
)

const (
	TravelCost     = 10
	FireCost       = 15
	ScanCost       = 5
	startingHull   = 100
	startingEnergy = 100
)

type HitFunc func(int, string)
type ScannedFunc func(string)

type Ship struct {
	hull        int `json:"shield"`
	energy      int `json:"energy"`
	hitFunc     HitFunc
	scannedFunc ScannedFunc
	name        string
	client      Client
}

func NewShip(client Client) *Ship {
	name := os.Getenv("SHIP_NAME")
	return &Ship{
		hull:   startingHull,
		energy: startingEnergy,
		name:   name,
		client: client,
	}
}

func NewShipFromJson(body io.ReadCloser) (*Ship, error) {
	decoder := json.NewDecoder(body)
	s := &Ship{}
	err := decoder.Decode(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Ship) SetHitFunc(hitFunc HitFunc) {
	s.hitFunc = hitFunc
}

func (s *Ship) SetScannedFunc(scannedFunc ScannedFunc) {
	s.scannedFunc = scannedFunc
}

func (s *Ship) GetHull() int {
	return s.hull
}

func (s *Ship) GetEnergy() int {
	return s.energy
}

func (s *Ship) CanTravel() bool {
	return s.energy > TravelCost
}

func (s *Ship) CanFire() bool {
	return s.energy > FireCost
}

func (s *Ship) CanScan() bool {
	return s.energy > ScanCost
}

func (s *Ship) ScanSector() ([]*Ship, []*Sector, error) {
	ships, sectors, newShip, err := s.client.ScanSector()
	if err != nil {
		s.Update(newShip)
	}
	return ships, sectors, err
}

func (s *Ship) Travel(sector *Sector) error {
	newShip, err := s.client.Travel(sector)
	if err != nil {
		return err
	}
	s.Update(newShip)
	return nil
}

func (s *Ship) Fire(target string) error {
	newShip, err := s.client.Fire(target)
	if err != nil {
		return err
	}
	s.Update(newShip)
	return nil
}

func (s *Ship) Update(newShip *Ship) {
	if newShip != nil {
		s.hull = newShip.GetHull()
		s.energy = newShip.GetEnergy()
	}
}

func (s *Ship) Hit(damage int, attacker string, newShip *Ship) {
	s.Update(newShip)
	if s.hitFunc != nil {
		s.hitFunc(damage, attacker)
	}
}

func (s *Ship) Scanned(attacker string) {
	if s.scannedFunc != nil {
		s.scannedFunc(attacker)
	}
}
