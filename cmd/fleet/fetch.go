package main

import (
	"github.com/TuhinNair/fleet/internal/data"
	"github.com/TuhinNair/fleet/internal/services/telematics"
)

func (a *application) fetchVehiclesData() error {
	serviceData, err := a.service.telematics.VehiclesSnapshot()
	if err != nil {
		return err
	}
	vehiclesDataModel := a.serviceToModelTransform(serviceData)

	for _, v := range vehiclesDataModel {
		err := a.model.Insert(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *application) serviceToModelTransform(d []*telematics.VehiclesData) []*data.Vehicle {
	var modelData []*data.Vehicle
	for _, v := range d {
		var vehicle data.Vehicle
		vehicle.ID = v.ID
		vehicle.Name = v.Name
		vehicle.Latitude = string(v.Latitude)
		vehicle.Longitude = string(v.Longitude)
		vehicle.AccurateAt = v.AccurateAt
		modelData = append(modelData, &vehicle)
	}
	return modelData
}
