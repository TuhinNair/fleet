package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/TuhinNair/fleet/internal/data"
	"github.com/TuhinNair/fleet/internal/jsonlog"
	"github.com/TuhinNair/fleet/internal/services/maps"
	"github.com/TuhinNair/fleet/internal/services/telematics"
	_ "github.com/lib/pq"
)

type application struct {
	model   data.VehicleModel
	service struct {
		telematics telematics.FleetDataFetcher
		maps       maps.Mapper
	}
	logger *jsonlog.Logger
}

func newApplication(cfg config) (*application, error) {
	var app application
	app.logger = jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app.logger.PrintInfo(fmt.Sprintf("db configuration: %+v", cfg.db), nil)
	db, err := configureDB(cfg.db)
	if err != nil {
		app.logger.PrintFatal(err, nil)
	}

	app.model = data.NewVehicleModel(db)
	app.service.telematics = telematics.NewFleetDataFetcher(
		cfg.services.telematics.baseURL,
		cfg.services.telematics.apiToken,
		cfg.services.telematics.pathSuffixes,
	)

	app.service.maps = maps.NewMapper(
		cfg.services.maps.baseURL,
		cfg.services.maps.apiToken,
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
