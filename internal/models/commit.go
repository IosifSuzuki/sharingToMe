package models

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

type CommitInfo struct {
	CreatedAt time.Time
	Messages []string
	Author string
}

func (c *CommitInfo)PrettyDate() string{
	return c.CreatedAt.Format("01/02/2006 03:04 PM")
}

func (c *CommitInfo)UnmarshalJSON(data []byte) error {
	var commitInfo struct{
		Commit struct {
			Author struct {
				FullName string `json:"name"`
				CreatedAt time.Time `json:"date"`
			}
			Message string `json:"message"`
		}
	}
	if err := json.Unmarshal(data, &commitInfo); err != nil {
		return err
	}
	c.Author = commitInfo.Commit.Author.FullName
	c.CreatedAt = commitInfo.Commit.Author.CreatedAt
	c.Messages = strings.Split(commitInfo.Commit.Message, "\n")
	var regExp = regexp.MustCompile(`[-;]`)
	for i, _ := range c.Messages {
		c.Messages[i] = strings.TrimSpace(regExp.ReplaceAllString(c.Messages[i], ""))
	}
	return nil
}
