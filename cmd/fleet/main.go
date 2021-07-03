package main

import "time"

func main() {
	cfg := loadConfig()

	app, err := newApplication(cfg)
	if err != nil {
		panic(err)
	}
	app.start(cfg.port)
}

func (a *application) start(port string) {
	a.poll(100 * time.Second)

	err := a.serve(port)
	if err != nil {
		a.logger.PrintFatal(err, nil)
	}
}
