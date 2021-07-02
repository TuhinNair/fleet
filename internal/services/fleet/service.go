package vehicles

import "encoding/json"

type FleetDataFetcher interface {
	VehiclesSnapshot() ([]*VehiclesData, error)
}

type VehiclesData struct {
	ID        string
	Name      string
	Latitude  json.Number
	Longitude json.Number
}
