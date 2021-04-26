package dbManager

import (
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/models"
	"errors"
	"net/url"
)

func insertPublisher(publisher *models.Publisher) error {
	stmt, err := DB.Prepare(`INSERT INTO "publisher" ("nickname", "email", "ip", "country_flag_url", "latitude", "longitude") VALUES($1, $2, $3, $4, $5, $6) RETURNING "id";`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(publisher.Nickname, publisher.Email, publisher.Ip, publisher.Flag.String(), publisher.Latitude, publisher.Longitude).Scan(&id)
	if err != nil {
		return err
	}
	publisher.Id = int(id)
	return err
}

func findPublisher(publisher *models.Publisher) error {
	stmt, err := DB.Prepare(`SELECT "id" FROM "publisher" WHERE "nickname" = $1 AND "ip" = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var publisherId = publisher.Id
	row := stmt.QueryRow(publisher.Nickname, publisher.Ip)
	_ = row.Scan(&publisherId)
	publisher.Id = publisherId
	return err
}

func insertPost(post models.Post) error {
	stmt, err := DB.Prepare(`INSERT INTO "post" ("publisher_id", "description", "file_path") VALUES($1, $2, $3);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post.Publisher.Id, post.Description, post.FilePath)
	return err
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
	stmtOfPost, err := DB.Prepare(`SELECT * FROM "post";`)
	if err != nil {
		return nil, err
	}
	stmtOfPublisher, err := DB.Prepare(`SELECT * FROM "publisher" WHERE "id" = $1`)
	if err != nil {
		return nil, err
	}


	defer stmtOfPost.Close()
	defer stmtOfPublisher.Close()
	rows, err := stmtOfPost.Query()
	if err != nil {
		return nil, err
	}
	var posts = make([]models.Post, 0)
	for rows.Next() {
		var (
			post models.Post
			publisherId int
		)
		err = rows.Scan(&post.Id, &publisherId, &post.Description, &post.FilePath)
		if err != nil {
			return posts, err
		}
		var publisher models.Publisher
		row := stmtOfPublisher.QueryRow(publisherId)
		var flagURL string
		err = row.Scan(
			&publisher.Id,
			&publisher.Nickname,
			&publisher.RegisteredAt,
			&publisher.Ip,
			&flagURL,
			&publisher.Latitude,
			&publisher.Longitude,
			&publisher.Email,
			)
		if err != nil {
			return posts, err
		}
		publisher.Flag, err = url.Parse(flagURL)
		if err != nil {
			return posts ,err
		}
		post.Publisher = &publisher
		posts = append(posts, post)
	}
	return posts, nil
}