package dtl

import (
	"encoding/json"
	"errors"
	"io"
)

type Action struct {
	Type    string
	Payload map[string]interface{}
	State   Ship
}

type HitPayload struct {
	Damage int
	Enemy  string
}

type ScanPayload struct {
	Enemy string
}

func NewActionFromJson(body io.ReadCloser) (*Action, error) {
	decoder := json.NewDecoder(body)
	a := &Action{}
	err := decoder.Decode(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Action) HitPayload() (HitPayload, error) {
	hp := HitPayload{}
	d := a.Payload["damage"]
	dInt, ok := d.(int)
	if !ok {
		return hp, errors.New("failed to parse hit payload")
	}
	e := a.Payload["enemy"]
	eStr, ok1 := e.(string)
	if !ok1 {
		return hp, errors.New("failed to parse hit payload")
	}
	hp.Damage = dInt
	hp.Enemy = eStr

	return hp, nil
}

func (a *Action) ScanPayload() (ScanPayload, error) {
	sp := ScanPayload{}
	e := a.Payload["enemy"]
	eStr, ok1 := e.(string)
	if !ok1 {
		return sp, errors.New("failed to parse scan payload")
	}
	sp.Enemy = eStr

	return sp, nil
}
