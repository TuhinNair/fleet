package main

import (
	"os"

	"github.com/TuhinNair/fleet/internal/services/telematics"
)

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type serviceConfig struct {
	telematics struct {
		baseURL      string
		apiToken     string
		pathSuffixes map[telematics.URLName]string
	}
	maps struct {
		baseURL  string
		apiToken string
	}
}
type config struct {
	db       dbConfig
	services serviceConfig
	port     string
}

func loadConfig() config {
	db := newDBConfig()
	services := newServicesConfig()
	port := ":" + os.Getenv("PORT")
	return config{db, services, port}
}

func newDBConfig() dbConfig {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("missing database URL in env")
	}
	return dbConfig{
		dsn:          dsn,
		maxOpenConns: 25,
		maxIdleConns: 25,
		maxIdleTime:  "15m",
	}
}

func newServicesConfig() serviceConfig {
	var scfg serviceConfig
	baseURL := os.Getenv("SAMSARA_BASE_URL")
	if baseURL == "" {
		panic("missing base api url")
	}
	pathSuffixes := make(map[telematics.URLName]string)
	pathSuffixes[telematics.VehiclesSnapshot] = os.Getenv("SAMSARA_VEHICLE_STATS_PATH_SUFFIX")
	if pathSuffixes[telematics.VehiclesSnapshot] == "" {
		panic("missing vehicle_snapshot path suffix")
	}
	apiToken := os.Getenv("SAMSARA_API_TOKEN")
	if apiToken == "" {
		panic("missing samsara API token")
	}
	scfg.telematics.baseURL = baseURL
	scfg.telematics.pathSuffixes = pathSuffixes
	scfg.telematics.apiToken = apiToken

	baseURL = os.Getenv("GOOGLE_MAPS_BASE_URL")
	if baseURL == "" {
		panic("missing base maps api url")
	}
	apiToken = os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiToken == "" {
		panic("missing maps API token")
	}
	scfg.maps.baseURL = baseURL
	scfg.maps.apiToken = apiToken
	return scfg
}
