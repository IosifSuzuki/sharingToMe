package main

import (
	"IosifSuzuki/sharingToMe/internal/dbManager"
	"IosifSuzuki/sharingToMe/internal/routing"
)

func main() {
	var router routing.Routing = routing.NewAPIRouter()
	defer dbManager.DB.Close()

	router.Setup()
	router.Run()
}
