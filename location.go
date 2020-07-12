package main

type Location struct {
	Longtitude float64
	Latitude   float64
}

func NewLocation(long, lat float64) *Location {
	return &Location{Longtitude: long, Latitude: lat}
}
