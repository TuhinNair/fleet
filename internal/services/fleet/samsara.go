package vehicles

var _ FleetDataFetcher = (*samsara)(nil)

type samsara struct {
}

func (s *samsara) VehiclesSnapshot() ([]*VehiclesData, error) {
	return nil, nil
}
