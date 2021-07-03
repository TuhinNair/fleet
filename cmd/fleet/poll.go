package main

import "time"

func (a *application) poll(every time.Duration) {
	run := time.Tick(every)
	go func() {
		for {
			select {
			case <-run:
				err := a.fetchVehiclesData()
				if err != nil {
					a.logger.PrintError(err, nil)
				}
			default:
				continue
			}
		}
	}()
}
