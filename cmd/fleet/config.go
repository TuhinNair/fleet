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
	baseURL      string
	apiToken     string
	pathSuffixes map[telematics.URLName]string
}
type config struct {
	db       dbConfig
	services serviceConfig
}

func loadConfig() config {
	db := newDBConfig()
	services := newServicesConfig()
	return config{db, services}
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
	return serviceConfig{baseURL, apiToken, pathSuffixes}
}
