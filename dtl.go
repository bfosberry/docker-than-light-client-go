package dtl

func New() (*Ship, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	s := NewShip(client)
	server := NewServer(s)
	server.Listen()
	return s, nil
}
