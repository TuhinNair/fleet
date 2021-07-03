package fleet

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type URLName string

type samsara struct {
	baseURL      string
	pathSuffixes map[URLName]string
	apiToken     string
	client       *http.Client
}

func newSamsaraService(baseURL, apiToken string, pathSuffixes map[URLName]string) *samsara {
	// the standard client uses cached TCP connections so encourages single client reuse
	client := &http.Client{}
	return &samsara{baseURL, pathSuffixes, apiToken, client}
}

func (s *samsara) endpoint(name URLName) string {
	suffix, ok := s.pathSuffixes[name]
	if !ok {
		// panic here as this is a development error
		panic("No path suffix exists for the given endpoint name: " + name)
	}
	return s.baseURL + suffix
}

func (s *samsara) do(req *http.Request) (*http.Response, error) {
	bearer := "Bearer " + s.apiToken
	req.Header.Add("Authorization", bearer)
	return s.client.Do(req)
}

func (s *samsara) getRequestWithTimeout(urlName URLName, timeout time.Duration) (*http.Request, context.CancelFunc, error) {
	url := s.endpoint(urlName)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, cancel, err
	}
	return req, cancel, nil
}

func (s *samsara) VehiclesSnapshot() ([]*VehiclesData, error) {
	req, cancel, err := s.getRequestWithTimeout("vehicles_snapshot", 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer cancel()

	resp, err := s.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return s.unmarshalVehiclesSnapshot(body)
}

func (s *samsara) unmarshalVehiclesSnapshot(data []byte) ([]*VehiclesData, error) {
	var expectedJSON = struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			GPS  struct {
				Time      time.Time   `json:"time"`
				Latitude  json.Number `json:"latitude"`
				Longitude json.Number `json:"longitude"`
			} `json:"gps"`
		} `json:"data"`
	}{}
	err := json.Unmarshal(data, &expectedJSON)
	if err != nil {
		return nil, err
	}

	var vehiclesData []*VehiclesData
	for _, vehicle := range expectedJSON.Data {
		var v VehiclesData
		v.ID = vehicle.ID
		v.Name = vehicle.Name
		v.AccurateAt = vehicle.GPS.Time
		v.Latitude = vehicle.GPS.Latitude
		v.Longitude = vehicle.GPS.Longitude
		vehiclesData = append(vehiclesData, &v)
	}
	return vehiclesData, nil
}
