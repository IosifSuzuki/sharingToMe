package defaults

import (
	"path/filepath"
	"regexp"
)

const (
	NewId = -1
	ContentType = "Content-Type"
)

var (
	BasePathToTemplates = filepath.Join("src", "templates", "*")
	RequestCompletedSuccessfully = "Request completed successfully"
	ConsumerAlreadyExist = "Such a consumer already exists"
	RequestFailed = "Request failed"
	UnauthorizedUser = "Your credentials are invalid"
	IsValidNickNamingFunction = regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString
	IsValidPasswordFunction = regexp.MustCompile(`^[0-9]{4}$`).MatchString
)
