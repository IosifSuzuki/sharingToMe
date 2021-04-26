package main

import "IosifSuzuki/sharingToMe/internal/routing"

func main() {
	var router routing.Routing = routing.NewWEBRouter()
	router.Setup()
	router.Run()
}
