package models

import (
	"net/url"
)

type Publisher struct {
	Id 			int
	Nickname 	string
	Email		string
	Ip 			string
	Flag 		*url.URL
	Latitude 	float32
	Longitude 	float32
}

