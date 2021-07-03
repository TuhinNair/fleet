package main

import "time"

func (a *application) poll(every time.Duration) {
	run := time.Tick(every)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-run:
	// 			a.logger.PrintInfo("polling fresh data", nil)
	// 			a.fetchVehiclesData()
	// 		default:
	// 			continue
	// 		}
	// 	}
	// }()
	for {
		select {
		case <-run:
			a.logger.PrintInfo("polling fresh data", nil)
			a.fetchVehiclesData()
		default:
			continue
		}
	}
}
