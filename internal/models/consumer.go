package models

import (
	"IosifSuzuki/sharingToMe/internal/defaults"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type Consumer struct {
	Id			int
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phoneNumber"`
	Password    string    `json:"password"`
	BirthDate   time.Time `json:"birthDate"`
	Reference   string	  `json:"reference,omitempty"`
	IpInfo 		IpInfo    `json:",omitempty"`
}

func (c *Consumer)ValidateForm() (bool, string) {
	if len(c.Username) == 0 {
		return false, "Username cannot be empty"
	} else if !defaults.IsValidNickNamingFunction(c.Username) {
		return false, "Username can contain letters, numbers and special a sign '-'"
	} else if len(c.PhoneNumber) > 13 || len(c.PhoneNumber) == 0 {
		return false, "Check your phone number, it contains an unsupported format"
	}
	if len(c.Password) == 0 {
		return false, "Password cannot be empty"
	} else if !defaults.IsValidPasswordFunction(c.Password) {
		return false, "Password must contain only 4 digits"
	}
	return true, ""
}

type Credential struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type Token struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phoneNumber"`
	Role        string `json:"role"`
	jwt.StandardClaims
}
