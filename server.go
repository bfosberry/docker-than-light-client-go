package dtl

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	DefaultPort = "8080"
)

type Server interface {
	Listen()
}

type server struct {
	port  string
	ship  *Ship
	token string
}

func NewServer(ship *Ship) (Server, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("No Token provided")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	return &server{
		port:  port,
		ship:  ship,
		token: token,
	}, nil
}

func (s *server) Listen() {
	http.HandleFunc("/_ping", s.ping)
	http.HandleFunc("/update", s.update)
	http.HandleFunc("/action", s.action)
	http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil)
}

func (s *server) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != s.token {
		w.WriteHeader(401)
		return
	}
	sh, err := NewShipFromJson(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	s.ship.Update(sh)
	w.WriteHeader(200)
}

func (s *server) action(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != s.token {
		w.WriteHeader(401)
		return
	}
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
			fmt.Printf("failed to to get payload: %s\n", err)
		} else {
			s.ship.Hit(hp.Damage, hp.Enemy, a.State)
		}
	case "scan":
		var sp ScanPayload
		sp, err = a.ScanPayload()
		if err != nil {
			fmt.Printf("failed to to get payload: %s\n", err)
		} else {
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
