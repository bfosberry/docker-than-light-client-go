package dtl

func New() (*Ship, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	s := NewShip(client)
	var server Server
	server, err = NewServer(s)
	if err != nil {
		return nil, err
	}
	server.Listen()
	return s, nil
}
