package models

import (
	"fmt"
	"mime/multipart"
	"regexp"
)

type NewPost struct {
	Nickname      string
	Email         string
	Description   string
	File          *multipart.File
	FileHeader    *multipart.FileHeader
	MaxSizeOfFile int64
}

func (n *NewPost) ValidatePostForm() map[string][]string {
	var errors = make(map[string][]string)
	var isEnglishLetters = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	errors["Nickname"] = make([]string, 0)
	if len(n.Nickname) == 0 {
		errors["Nickname"] = append(errors["Nickname"], "Nickname is required field")
	}
	if len(n.Nickname) > 32 {
		errors["Nickname"] = append(errors["Nickname"], "Nickname must contains less then or equal 32 letters")
	}
	if !isEnglishLetters(n.Nickname) {
		errors["Nickname"] = append(errors["Nickname"], "Nickname must contains only english letters")
	}
	errors["Email"] = make([]string, 0)
	var isEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString
	if len(n.Email) != 0 && !isEmail(n.Email) {
		errors["Email"] = append(errors["Email"], "Not valid email")
	}
	errors["File"] = make([]string, 0)
	if n.File == nil || n.FileHeader == nil {
		errors["File"] = append(errors["File"], "You must upload audio file")
	}
	if n.FileHeader != nil && n.FileHeader.Size > n.MaxSizeOfFile {
		errors["File"] = append(errors["File"], fmt.Sprintf("Your file has %d bytes, but your file must contains less them %d bytes", n.FileHeader.Size, n.MaxSizeOfFile))
	}
	if len(errors["Nickname"]) == 0 {
		delete(errors, "Nickname")
	}
	if len(errors["File"]) == 0 {
		delete(errors, "File")
	}
	if len(errors["Email"]) == 0 {
		delete(errors, "Email")
	}
	return errors
}
