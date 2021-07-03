package main

import "time"

func (a *application) poll(every time.Duration) {
	run := time.Tick(every)
	go func() {
		for {
			select {
			case <-run:
				a.logger.PrintInfo("polling fresh data", nil)
				err := a.fetchVehiclesData()
				a.logger.PrintError(err, nil)
			default:
				continue
			}
		}
	}()
}
