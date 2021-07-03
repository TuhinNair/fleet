package data

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestVehiclesDataSource(t *testing.T) {
	testDB := setupTestDB(t)
	model := NewVehicleModel(testDB)
	tests := map[string]func(t *testing.T){
		"insert data successfully": func(t *testing.T) {
			data := dummyVehiclesData()
			for _, vehicle := range data {
				err := model.Insert(vehicle)
				if err != nil {
					t.Fatal(err)
				}
			}
		},
		"retrieve all data successfully": func(t *testing.T) {
			data := dummyVehiclesData()
			for _, vehicle := range data {
				err := model.Insert(vehicle)
				if err != nil {
					t.Fatal(err)
				}
			}
			result, err := model.GetAll()
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range result {
				if v.ID == "1" {
					if !compareVehicleFields(v, data[0]) {
						t.Fatal("GetAll result data mismatch for vehicle with id:" + v.ID)
					}
				}
				if v.ID == "2" {
					if !compareVehicleFields(v, data[1]) {
						t.Fatal("GetAll result data mismatch for vehicle with id:" + v.ID)
					}
				}
				if v.ID == "3" {
					if !compareVehicleFields(v, data[2]) {
						t.Fatal("GetAll result data mismatch for vehicle with id:" + v.ID)
					}
				}
			}
		},
		"update vehicle data is successfull and correct": func(t *testing.T) {
			data := dummyVehiclesData()
			for _, vehicle := range data {
				err := model.Insert(vehicle)
				if err != nil {
					t.Fatal(err)
				}
			}

			updatedVehicle := data[0]
			updatedVehicle.Latitude = "99.00"

			err := model.Insert(updatedVehicle)
			if err != nil {
				t.Fatal(err)
			}

			rows, err := testDB.Query(`SELECT id, latitude FROM vehicles WHERE id = '1'`)
			if err != nil {
				t.Fatal(err)
			}

			numRows := 0
			if rows.Next() {
				numRows++
				var res Vehicle
				rows.Scan(&res.ID, &res.Latitude)
				if res.ID != "1" {
					t.Fatal("unexpected vehicle id after update")
				}
				if res.Latitude != "99.00" {
					t.Fatal("unexpected vehicle latitude after update")
				}
			}
			if numRows != 1 {
				t.Fatalf("update was unsuccessfult. Unexpected number of result rows. Got=%d,Expected=%d", numRows, 1)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}

func compareVehicleFields(a, b *Vehicle) bool {
	return a.ID == b.ID && a.Name == b.Name && a.Latitude == b.Latitude && a.Longitude == b.Longitude
}

func dummyVehiclesData() []*Vehicle {
	return []*Vehicle{
		{
			ID:         "1",
			Name:       "vehicle-1",
			Latitude:   "1.00",
			Longitude:  "2.00",
			AccurateAt: time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "2",
			Name:       "vehicle-2",
			Latitude:   "10.00",
			Longitude:  "20.00",
			AccurateAt: time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "3",
			Name:       "vehicle-3",
			Latitude:   "100.00",
			Longitude:  "200.00",
			AccurateAt: time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
}

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("missing test database url")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
