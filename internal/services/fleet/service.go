package fleet

import (
	"encoding/json"
	"time"
)

type FleetDataFetcher interface {
	VehiclesSnapshot() ([]*VehiclesData, error)
}

func NewFleetDataFetcher(baseURL, apiToken string, pathSuffixes map[URLName]string) FleetDataFetcher {
	return newSamsaraService(baseURL, apiToken, pathSuffixes)
}

type VehiclesData struct {
	ID         string
	Name       string
	Latitude   json.Number
	Longitude  json.Number
	AccurateAt time.Time
}
