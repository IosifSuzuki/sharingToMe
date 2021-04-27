package dbManager

import (
	"IosifSuzuki/sharingToMe/internal/models"
	"net/url"
)

func getAllPosts() ([]models.Post, error) {
	stmtOfPosts, err := DB.Prepare(`SELECT * FROM "post"`)
	if err != nil {
		return nil, err
	}
	defer stmtOfPosts.Close()
	rows, err := stmtOfPosts.Query()
	if err != nil {
		return nil, err
	}
	var posts = make([]models.Post, 0)
	for rows.Next() {
		var (
			post models.Post
		)
		post.Publisher = new(models.Publisher)
		err = rows.Scan(
			&post.Id,
			&post.Publisher.Id,
			&post.Description,
			&post.FilePath,
			&post.CreatedAt,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func addPublisherToPost(posts []models.Post ) error {
	stmtOfPublisher, err := DB.Prepare(`SELECT * FROM "publisher" WHERE "id" = $1`)
	if err != nil {
		return err
	}
	defer stmtOfPublisher.Close()
	for i := range posts {
		var (
			row = stmtOfPublisher.QueryRow(posts[i].Publisher.Id)
			flagURL string
		)
		err = row.Scan(
			&posts[i].Publisher.Id,
			&posts[i].Publisher.Nickname,
			&posts[i].Publisher.Ip,
			&flagURL,
			&posts[i].Publisher.Latitude,
			&posts[i].Publisher.Longitude,
			&posts[i].Publisher.Email,
		)
		posts[i].Publisher.Flag, err = url.Parse(flagURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertPublisher(publisher *models.Publisher) error {
	stmt, err := DB.Prepare(`INSERT INTO "publisher" (
"nickname", "email", "ip", "country_flag_url", "latitude", "longitude"
) VALUES($1, $2, $3, $4, $5, $6) RETURNING "id"`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(
		publisher.Nickname,
		publisher.Email,
		publisher.Ip,
		publisher.Flag.String(),
		publisher.Latitude,
		publisher.Longitude,
		).Scan(&id)
	if err != nil {
		return err
	}
	publisher.Id = id
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
	stmt, err := DB.Prepare(`INSERT INTO "post" (
"publisher_id", "description", "file_path"
) VALUES($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post.Publisher.Id, post.Description, post.FilePath)
	return err
}