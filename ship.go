package dtl

import (
	"os"
)

type Ship interface {
	CanTravel() bool
	CanFire() bool
	CanScan() bool
	ScanSector() ([]Ship, []Sector, error)
	Travel(Sector) error
	Fire(string) error
	Update(Ship)
	Hit(int, string, Ship)
	Scanned(string)
	GetHull() int
	GetEnergy() int
}

const (
	TravelCost     = 10
	FireCost       = 15
	ScanCost       = 5
	startingHull   = 100
	startingEnergy = 100
)

type HitFunc func(int, string)
type ScannedFunc func(string)

type ship struct {
	hull        int
	energy      int
	hitFunc     HitFunc
	scannedFunc ScannedFunc
	name        string
	client      Client
}

func NewShip(client Client) Ship {
	name := os.Getenv("SHIP_NAME")
	return &ship{
		hull:   startingHull,
		energy: startingEnergy,
		name:   name,
		client: client,
	}
}

func (s *ship) SetHitFunc(hitFunc HitFunc) {
	s.hitFunc = hitFunc
}

func (s *ship) SetScannedFunc(scannedFunc ScannedFunc) {
	s.scannedFunc = scannedFunc
}

func (s *ship) GetHull() int {
	return s.hull
}

func (s *ship) GetEnergy() int {
	return s.energy
}

func (s *ship) CanTravel() bool {
	return s.energy > TravelCost
}

func (s *ship) CanFire() bool {
	return s.energy > FireCost
}

func (s *ship) CanScan() bool {
	return s.energy > ScanCost
}

func (s *ship) ScanSector() ([]Ship, []Sector, error) {
	ships, sectors, newShip, err := s.client.ScanSector()
	if err != nil {
		s.Update(newShip)
	}
	return ships, sectors, err
}

func (s *ship) Travel(sector Sector) error {
	newShip, err := s.client.Travel(sector)
	if err != nil {
		return err
	}
	s.Update(newShip)
	return nil
}

func (s *ship) Fire(target string) error {
	newShip, err := s.client.Fire(target)
	if err != nil {
		return err
	}
	s.Update(newShip)
	return nil
}

func (s *ship) Update(newShip Ship) {
	s.hull = newShip.GetHull()
	s.energy = newShip.GetEnergy()
}

func (s *ship) Hit(damage int, attacker string, newShip Ship) {
	s.Update(newShip)
	s.hitFunc(damage, attacker)
}

func (s *ship) Scanned(attacker string) {
	s.scannedFunc(attacker)
}
