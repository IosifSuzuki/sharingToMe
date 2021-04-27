package models

import (
	"fmt"
	"path/filepath"
	"time"
)

//CREATE TABLE post (
//id SERIAL PRIMARY KEY,
//post_id INT NOT NULL,
//description TEXT,
//file_path TEXT NOT NULL
//);

type Post struct {
	Id 			int
	Description string
	FilePath 	string
	Publisher 	*Publisher
	CreatedAt	time.Time
}

func (p *Post)RemoteURL() string {
	return fmt.Sprintf("/static/files/%s", filepath.Base(p.FilePath))
}

func (p *Post)PrettyDate() string {
	return p.CreatedAt.Format("01/02/2006 03:04 PM")
}
