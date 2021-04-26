package fileManager

import (
	"crypto/md5"
	"crypto/rand"
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
	var byte = make([]byte, 8)
	rand.Read(byte)
	var (
		fileName = fmt.Sprintf("%X", md5.New().Sum(byte))
		fullPath = filepath.Join(BaseDir, "src", "assets", "files", fmt.Sprintf("%s%s", fileName, extension))
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
