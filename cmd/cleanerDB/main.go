package main

import (
	"IosifSuzuki/sharingToMe/internal/dbManager"
	"IosifSuzuki/sharingToMe/internal/fileManager"
)

func main() {
	filePaths, err := dbManager.ClearOldData()
	if err != nil {
		panic(err)

	}
	for _, filePath := range filePaths {
		err = fileManager.RemoveFile(filePath)
		if err != nil {
			panic(err)
		}
	}
}
