package dtl

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

const (
	TravelCost     = 15
	FireCost       = 12
	ScanCost       = 5
	startingShield = 100
	startingEnergy = 100
	rechargeRate   = 10
	maxEnergy      = 100
)

type HitFunc func(int, string)
type ScannedFunc func(string)

type Ship struct {
	Shield      int `json:"shield"`
	Energy      int `json:"energy"`
	hitFunc     HitFunc
	scannedFunc ScannedFunc
	Name        string `json:"name"`
	client      Client
}

func NewShip(client Client) *Ship {
	name := os.Getenv("SHIP_NAME")
	s := &Ship{
		Shield: startingShield,
		Energy: startingEnergy,
		Name:   name,
		client: client,
	}
	go s.recharge()
	return s
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

func (s *Ship) recharge() {
	for {
		time.Sleep(1 * time.Second)
		s.Energy = s.Energy + rechargeRate
		if s.Energy > maxEnergy {
			s.Energy = maxEnergy
		}
	}
}

func (s *Ship) SetHitFunc(hitFunc HitFunc) {
	s.hitFunc = hitFunc
}

func (s *Ship) SetScannedFunc(scannedFunc ScannedFunc) {
	s.scannedFunc = scannedFunc
}

func (s *Ship) CanTravel() bool {
	return s.Energy > TravelCost
}

func (s *Ship) CanFire() bool {
	return s.Energy > FireCost
}

func (s *Ship) CanScan() bool {
	return s.Energy > ScanCost
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
		s.Shield = newShip.Shield
		s.Energy = newShip.Energy
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
