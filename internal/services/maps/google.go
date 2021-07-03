package maps

type googleMaps struct {
	baseURL  string
	apiToken string
}

func newGoogleMapsService(baseURL, apiToken string) Mapper {
	return &googleMaps{baseURL, apiToken}
}

func (g *googleMaps) MapURL(markers []Marker) string {
	return g.baseURL + "?size=500x500&key=" + g.apiToken + g.buildMarkerParams(markers)
}

func (g *googleMaps) buildMarkerParams(markers []Marker) string {
	var markerParams string
	for _, m := range markers {
		markerParams += "&markers=label:" + m.Label + "%7C" + m.Location
	}
	return markerParams
}
