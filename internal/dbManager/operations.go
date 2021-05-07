package dbManager

import (
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/models"
	"IosifSuzuki/sharingToMe/internal/utility"
	"database/sql"
	"errors"
	"time"
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

func addPublisherToPost(posts []models.Post) error {
	stmtOfPublisher, err := DB.Prepare(`SELECT * FROM "publisher" WHERE "id" = $1`)
	if err != nil {
		return err
	}
	defer stmtOfPublisher.Close()
	for i := range posts {
		var (
			row     = stmtOfPublisher.QueryRow(posts[i].Publisher.Id)
		)
		err = row.Scan(
			&posts[i].Publisher.Id,
			&posts[i].Publisher.Nickname,
			&posts[i].Publisher.IpInfo.IP,
			&posts[i].Publisher.IpInfo.CountryFlag,
			&posts[i].Publisher.IpInfo.Latitude,
			&posts[i].Publisher.IpInfo.Longitude,
			&posts[i].Publisher.Email,
		)
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
		publisher.IpInfo.IP,
		publisher.IpInfo.CountryFlag,
		publisher.IpInfo.Latitude,
		publisher.IpInfo.Longitude,
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
	row := stmt.QueryRow(publisher.Nickname, publisher.IpInfo.IP)
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

func isExistPublisher(publisher models.Publisher) (bool, error) {
	stmt, err := DB.Prepare(`SELECT "id" FROM "publisher" WHERE "nickname" = $1 AND "ip" = $2`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var publisherId = publisher.Id
	row := stmt.QueryRow(publisher.Nickname, publisher.IpInfo.IP)
	_ = row.Scan(&publisherId)
	return publisherId != defaults.NewId, err
}

func clearOldData() ([]string, error) {
	var (
		dateForDeletePost = time.Now().Add(- 3 * 24 * time.Hour)
		files             = make([]string, 0, 0)
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

func allowCreatePost(ip string) (bool, error) {
	var (
		year, month, day = time.Now().Date()
		todayBeginDay    = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
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

func isExistConsumer(consumer models.Consumer) (bool, error) {
	stmt, err := DB.Prepare(`SELECT id FROM "consumer" WHERE "phone_number" = $1`)
	if err != nil {
		return false, nil
	}
	defer stmt.Close()
	var row = stmt.QueryRow(consumer.PhoneNumber)
	err = row.Scan(&consumer.Id)
	if err != nil && err == sql.ErrNoRows {
		return false, nil
	}
	return consumer.Id != defaults.NewId, err
}

func saveConsumer(consumer models.Consumer) error {
	stmt, err := DB.Prepare(`INSERT INTO "consumer"("nickname", "phone_number", "password_hash", 
    "birth_date", "reference", "ip", "country_flag_url", "latitude", "longitude") 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`)
	if err != nil {
		return nil
	}
	defer stmt.Close()
	_, err = stmt.Exec(consumer.Username, consumer.PhoneNumber, consumer.Password, consumer.BirthDate,
		consumer.Reference, consumer.IpInfo.IP, consumer.IpInfo.CountryFlag, consumer.IpInfo.Latitude, consumer.IpInfo.Longitude)
	return err
}


func fetchConsumer(credential models.Credential) (*models.Consumer, error) {
	stmt, err := DB.Prepare(`SELECT * FROM "consumer" WHERE "phone_number" = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var (
		row = stmt.QueryRow(credential.PhoneNumber)
		consumer models.Consumer
	)
	err = row.Scan(&consumer.Id, &consumer.Username, &consumer.PhoneNumber, &consumer.Password,
		&consumer.BirthDate, &consumer.Reference, &consumer.IpInfo.IP, &consumer.IpInfo.CountryFlag,
		&consumer.IpInfo.Latitude, &consumer.IpInfo.Longitude)
	if err == sql.ErrNoRows {
		return nil, errors.New("User with the phone number does not exist")
	}
	var isSignIn = utility.CompareHashPassword(credential.Password, consumer.Password)
	if isSignIn {
		return &consumer, err
	} else {
		return nil, errors.New("Your password does not match with the phone number")
	}
}
