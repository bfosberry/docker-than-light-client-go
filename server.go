package dtl

import (
	"fmt"
	"net/http"
)

const (
	DefaultPort = 80
)

type Server interface {
	Listen()
}

type server struct {
	port int
	ship Ship
}

func NewServer(ship Ship) Server {
	return &server{
		port: DefaultPort,
		ship: ship,
	}
}

func (s *server) Listen() {
	http.HandleFunc("/_ping", s.ping)
	http.HandleFunc("/update", s.update)
	http.HandleFunc("/action", s.action)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *server) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	sh, err := NewShipFromJson(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	s.ship.Update(sh)
	w.WriteHeader(200)
}

func (s *server) action(w http.ResponseWriter, r *http.Request) {
	a, err := NewActionFromJson(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	switch a.Type {
	case "hit":
		var hp HitPayload
		hp, err = a.HitPayload()
		if err != nil {
			s.ship.Hit(hp.Damage, hp.Enemy, a.State)
		}
	case "scan":
		var sp ScanPayload
		sp, err = a.ScanPayload()
		if err != nil {
			s.ship.Scanned(sp.Enemy)
			s.ship.Update(a.State)
		}
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
