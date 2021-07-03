package data

import (
	"context"
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

func (m VehicleModel) GetAll() ([]*Vehicle, error) {
	query := `
		SELECT id, name, latitude, longitude, accurate_at, updated_at
		FROM vehicles`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Vehicle
	for rows.Next() {
		var vehicle Vehicle
		err := rows.Scan(
			&vehicle.ID,
			&vehicle.Name,
			&vehicle.Latitude,
			&vehicle.Longitude,
			&vehicle.AccurateAt,
			&vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &vehicle)
	}
	return result, nil
}

func (m VehicleModel) Insert(v *Vehicle) error {
	query := `
		INSERT INTO vehicles (id, name, latitude, longitude, accurate_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		latitude = EXCLUDED.latitude,
		longitude = EXCLUDED.longitude,
		accurate_at = EXCLUDED.accurate_at,
		updated_at = EXCLUDED.updated_at`

	v.UpdatedAt = time.Now()
	args := []interface{}{
		v.ID,
		v.Name,
		v.Latitude,
		v.Longitude,
		v.AccurateAt,
		v.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}
