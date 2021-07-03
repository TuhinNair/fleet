package main

import "time"

func main() {
	cfg := loadConfig()
	app, err := newApplication(cfg)
	if err != nil {
		panic(err)
	}
	app.start()
}

func (a *application) start() {
	a.poll(10 * time.Second)
}
