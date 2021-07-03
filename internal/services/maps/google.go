package maps

type googleMaps struct {
	baseURL  string
	apiToken string
}

func newGoogleMapsService(baseURL, apiToken string) Mapper {
	return &googleMaps{baseURL, apiToken}
}

func (g *googleMaps) MapURL(markers []Marker) string {
	return g.baseURL + "?key=" + g.apiToken + g.buildMarkerParams(markers)
}

func (g *googleMaps) buildMarkerParams(markers []Marker) string {
	var markerParams string
	for _, m := range markers {
		markerParams += "&markers=label:" + m.Label + "|" + m.Location
	}
	return markerParams
}
