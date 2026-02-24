package model

import "fmt"

type Server struct {
	ID      string  `graphql:"_id"`
	Country *string `graphql:"country"`
	City    *string `graphql:"city"`
	IP      string  `graphql:"ip"`
	Name    string  `graphql:"name"`
}

func (s Server) GetID() string {
	return "server-" + s.ID
}

func (s Server) String() string {
	var identifier string

	if s.Country != nil {
		identifier = *s.Country
		if s.City != nil {
			identifier = fmt.Sprintf("%s, %s", *s.City, *s.Country)
		}
	} else {
		identifier = s.IP
	}

	return fmt.Sprintf("%s (%s)", s.Name, identifier)
}

func (s Server) IsAvailable() bool {
	return true
}

type Servers []Server

func (s Servers) Header() []string {
	return []string{"ID", "Name", "Location", "IP"}
}

func (s Servers) Rows() [][]string {
	rows := make([][]string, len(s))
	for i, server := range s {
		location := server.IP
		if server.Country != nil {
			location = *server.Country
			if server.City != nil {
				location = fmt.Sprintf("%s, %s", *server.City, *server.Country)
			}
		}
		rows[i] = []string{server.GetID(), server.Name, location, server.IP}
	}
	return rows
}
