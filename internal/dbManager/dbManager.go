package dbManager
import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB = makeConnectionToDB()

func makeConnectionToDB() *sql.DB {
	var dbInfo = configuration.Configuration.MainDB
	var psqlConnectionText = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",  dbInfo.Host, dbInfo.Port, dbInfo.Username, dbInfo.Password, dbInfo.DBName)
	db, err := sql.Open("postgres", psqlConnectionText)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}
