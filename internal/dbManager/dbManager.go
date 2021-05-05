package dbManager

import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/models"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
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
	stmt, err := DB.Prepare(`SELECT "id" FROM "publisher" WHERE "nickname" = $1 AND "ip" = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var publisherId = publisher.Id
	row := stmt.QueryRow(publisher.Nickname, publisher.Ip)
	_ = row.Scan(&publisherId)
	return publisherId != defaults.NewId, err
}

func ClearOldData() ([]string, error) {
	var (
		dateForDeletePost = time.Now().Add(- 3 * 24 * time.Hour)
		files = make([]string, 0, 0)
	)
	stmt, err := DB.Prepare(`DELETE FROM "post" WHERE "created_at" < $1 RETURNING "file_path"`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(dateForDeletePost)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var filePath string
		err = rows.Scan(&filePath)
		if err != nil {
			return files, err
		}
		files = append(files, filePath)
	}
	return files, nil
}

func AllowCreatePost(ip string) (bool, error) {
	var (
		year, month, day = time.Now().Date()
		todayBeginDay = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	)
	stmt, err := DB.Prepare(`SELECT COUNT(*) FROM 
    "post" INNER JOIN "publisher" ON post.publisher_id = publisher.id
	WHERE "created_at" > $1 AND "ip" = $2`)
	if err != nil {
		return false, nil
	}
	defer stmt.Close()
	var row = stmt.QueryRow(todayBeginDay, ip)
	var countOfPost int
	err = row.Scan(&countOfPost)
	if err != nil {
		return false, nil
	}
	return countOfPost < 3, nil
}