package dtl

func New() (Ship, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return NewShip(client), nil
}
