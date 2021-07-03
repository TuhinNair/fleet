package maps

type Marker struct {
	Label    string
	Location string
}

type Mapper interface {
	MapURL(markers []Marker) string
}

func NewMapper(baseURL, apiToken string) Mapper {
	return newGoogleMapsService(baseURL, apiToken)
}
