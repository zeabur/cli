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
