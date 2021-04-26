package models

import (
	"fmt"
	"path/filepath"
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
}

func (p *Post)RemoteURL() string {
	return fmt.Sprintf("/static/files/%s", filepath.Base(p.FilePath))
}
