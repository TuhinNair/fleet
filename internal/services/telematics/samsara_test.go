package telematics

import (
	"os"
	"testing"
)

func TestSamsaraJSONResponseDecoding(t *testing.T) {
	service := newSamaraTestService(t)

	tests := map[string]func(st *testing.T){
		"vehicle_snapshot response successfully decodes": func(t *testing.T) {
			vehicles, err := service.VehiclesSnapshot()
			if err != nil {
				t.Fatal(err)
			}

			for _, v := range vehicles {
				if v.ID == "" {
					t.Fatal("vehicle ID is empty")
				}
				if v.Name == "" {
					t.Fatal("vehicle name is empty")
				}
				if v.Latitude == "" {
					t.Fatal("vehicle latitude is empty")
				}
				if v.Longitude == "" {
					t.Fatal("vehicle longitude is empty")
				}
				if v.AccurateAt.IsZero() {
					t.Fatal("vehicle data timestamp is zero")
				}
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}

func newSamaraTestService(t *testing.T) *samsara {
	/* Although I used the same URLs for testing here,
	I used different env vars to emphasize this is to be
	against a sandbox */
	baseURL := os.Getenv("SAMSARA_TEST_BASE_URL")
	if baseURL == "" {
		t.Fatal("Missing samsara base url")
	}
	urlSuffixes := make(map[URLName]string)
	urlSuffixes[VehiclesSnapshot] = os.Getenv("SAMSARA_VEHICLE_STATS_PATH_SUFFIX")
	if urlSuffixes[VehiclesSnapshot] == "" {
		t.Fatal("missing vehicle_snapshot path suffix")
	}
	apiToken := os.Getenv("SAMSARA_TEST_API_TOKEN")
	if apiToken == "" {
		t.Fatal("missing samsara API token")
	}
	return newSamsaraService(baseURL, apiToken, urlSuffixes)
}
