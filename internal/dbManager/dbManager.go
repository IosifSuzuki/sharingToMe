package dbManager

import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/models"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var DB = makeConnectionToDB()

func makeConnectionToDB() *sql.DB {
	var dbInfo = configuration.Configuration.MainDB
	var psqlConnectionText = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.Host,
		dbInfo.Port,
		dbInfo.Username,
		dbInfo.Password,
		dbInfo.DBName,
	)
	db, err := sql.Open("postgres", psqlConnectionText)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func WritePostToDB(post models.Post) error {
	if post.Publisher == nil {
		return errors.New("a publisher of post can't be a nil")
	}
	if err := findPublisher(post.Publisher); err != nil {
		return err
	}
	if post.Publisher.Id == defaults.NewId {
		var err = insertPublisher(post.Publisher)
		if err != nil {
			return err
		}
	}
	var err = insertPost(post)
	if err != nil {
		return err
	}
	return nil
}

func ReadPosts() ([]models.Post, error) {
	posts, err := getAllPosts()
	if err != nil {
		return nil, err
	}
	err = addPublisherToPost(posts)
	return posts, err
}

func IsExistPublisher(publisher models.Publisher) (bool, error) {
	return isExistPublisher(publisher)
}

func ClearOldData() ([]string, error) {
	return clearOldData()
}

func AllowCreatePost(ip string) (bool, error) {
	return allowCreatePost(ip)
}

func IsExistConsumer(consumer models.Consumer) (bool, error) {
	return isExistConsumer(consumer)
}

func SaveConsumer(consumer models.Consumer) error {
	return saveConsumer(consumer)
}

func FetchConsumer(credential models.Credential) (*models.Consumer, error) {
	return fetchConsumer(credential)
}

