package data

import (
	"database/sql"
	"time"
)

type Vehicle struct {
	ID         string
	Name       string
	Latitude   string
	Longitude  string
	AccurateAt time.Time
	UpdatedAt  time.Time
}

type VehicleModel struct {
	DB *sql.DB
}
