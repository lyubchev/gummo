package main

type Location struct {
	Longtitude float64 `json:"longtitude,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
}

func NewLocation(long, lat float64) *Location {
	return &Location{Longtitude: long, Latitude: lat}
}
