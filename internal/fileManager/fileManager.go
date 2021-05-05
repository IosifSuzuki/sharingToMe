package fileManager

import (
	"IosifSuzuki/sharingToMe/internal/utility"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var BaseDir = getBaseDir()

func getBaseDir() string {
	baseDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return baseDir
}

func SaveMediaFile(file multipart.File, extension string) (*string, error) {
	var (
		fileName, _ = utility.NewUUID()
		fullPath = filepath.Join(BaseDir, "src", "files", fmt.Sprintf("%s%s", fileName, extension))
	)
	destinationFile, err := os.Create(fullPath)
	defer destinationFile.Close()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(destinationFile, file); err != nil {
		return nil, err
	}

	return &fullPath, nil
}

func RemoveFile(filePath string) error {
	var err = os.Remove(filePath)
	return err
}
