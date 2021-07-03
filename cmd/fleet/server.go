package main

import (
	"fmt"
	"net/http"

	"github.com/TuhinNair/fleet/internal/data"
	"github.com/TuhinNair/fleet/internal/services/maps"
)

func (a *application) serve(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.mapHandler)

	a.logger.PrintInfo("starting server on port"+port, nil)
	return http.ListenAndServe(port, a.globalRateLimit(mux))
}

func (a *application) mapHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := a.model.GetAll()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var markers []maps.Marker
	for _, v := range vehicles {
		markers = append(markers, newVehicleMapMarker(v))
	}
	mapURL := a.service.maps.MapURL(markers)

	htmlImg := `<img src="` + mapURL + `">`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = fmt.Fprint(w, htmlImg)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func newVehicleMapMarker(v *data.Vehicle) maps.Marker {
	return maps.Marker{
		Label:    v.ID,
		Location: v.Latitude + "," + v.Longitude,
	}
}
