package models

type App struct {
	SecretKey string
}

type Server struct {
	Host string
	Port int
}

type DB struct {
	Host 		string
	Port 		int
	Username 	string
	Password 	string
	DBName 		string
}

type ConfigurationFile struct {
	WebServer 	Server
	ApiServer 	Server
	MainDB 		DB
	AppInfo 	App
}
