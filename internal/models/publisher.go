package models

import (
	"net/url"
	"time"
)

type Publisher struct {
	Id int
	Nickname string
	Email string
	RegisteredAt *time.Time
	Ip string
	Flag *url.URL
	Latitude float32
	Longitude float32
}

func (p *Publisher)PrettyDate() string {
	return p.RegisteredAt.Format("01/02/2006 03:04 PM")
}

