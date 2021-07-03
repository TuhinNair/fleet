package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/TuhinNair/fleet/internal/data"
	"github.com/TuhinNair/fleet/internal/jsonlog"
	"github.com/TuhinNair/fleet/internal/services/telematics"
	_ "github.com/lib/pq"
)

type application struct {
	model   data.VehicleModel
	service telematics.FleetDataFetcher
	logger  *jsonlog.Logger
}

func newApplication(cfg config) (*application, error) {
	var app application
	app.logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := configureDB(cfg.db)
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

	app.model = data.NewVehicleModel(db)
	app.service = telematics.NewFleetDataFetcher(
		cfg.services.baseURL,
		cfg.services.apiToken,
		cfg.services.pathSuffixes,
	)

	err = app.fetchVehiclesData()
	if err != nil {
		app.logger.PrintError(err, nil)
		return nil, err
	}
	return &app, nil
}

func configureDB(cfg dbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)

	duration, err := time.ParseDuration(cfg.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, err
}
