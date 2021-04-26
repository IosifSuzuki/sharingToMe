package routing

import (
	"IosifSuzuki/sharingToMe/internal/dbManager"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/fileManager"
	"IosifSuzuki/sharingToMe/internal/models"
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"encoding/json"
	template2 "html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	contentType = "Content-Type"
)

var (
	basePathToTemplates = filepath.Join("src", "templates")
)

// API Handlers

func indexGetHandler(w http.ResponseWriter, _ *http.Request) {
	var notification = models.Notification{
		Title: "Welcome to us",
		Description: "This service provide send & store date to server",
	}
	w.Header().Set(contentType, "application/json")
	err := json.NewEncoder(w).Encode(notification)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// WEB

func homeGetHandler(w http.ResponseWriter, r *http.Request) {
	baseDir, err := os.Getwd()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := template2.Must(template2.ParseFiles(filepath.Join(baseDir, basePathToTemplates, "home.html")))
	err = tmpl.Execute(w, nil)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func homePostHandler(w http.ResponseWriter, r *http.Request) {
	const (
		nickNameKey = "nickname"
		emailKey = "email"
		descriptionKey = "description"
		fileKey = "file"
	)
	var err = r.ParseForm()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = r.ParseMultipartForm(5 * 1024 * 1024)
	var (
		nicknameValue = r.FormValue(nickNameKey)
		emailValue = r.FormValue(emailKey)
		descriptionValue = r.FormValue(descriptionKey)
	)
	file, fileHandler, err := r.FormFile(fileKey)
	if len(nicknameValue) == 0 || fileHandler == nil || err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	var fileExtension = filepath.Ext(fileHandler.Filename)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	filepath, err := fileManager.SaveMediaFile(file, fileExtension)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var flagURL, _ = url.Parse("https://assets.ipstack.com/flags/ua.svg")
	publisherId, err := utility.GetIP(r)
	if err != nil {
		publisherId = "0.0.0.0"
	}
	var publisher = models.Publisher{
		Id: defaults.NewId,
		Nickname: nicknameValue,
		Email: emailValue,
		RegisteredAt: nil,
		Ip: publisherId,
		Flag: flagURL,
		Latitude: 48.62828063964844,
		Longitude: 22.514659881591797,
	}
	var post = models.Post{
		Id: defaults.NewId,
		Description: descriptionValue,
		FilePath: *filepath,
		Publisher: &publisher,
	}
	if err := dbManager.WritePostToDB(post); err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sources", http.StatusSeeOther)
}

func listOfSourcesGetHandler(w http.ResponseWriter, r *http.Request) {
	baseDir, err := os.Getwd()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := template2.Must(template2.ParseFiles(filepath.Join(baseDir, basePathToTemplates, "listOfSources.html")))
	posts, err := dbManager.ReadPosts()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, struct {
		Title string
		Posts []models.Post
	} {
		Title: "List of Sources",
		Posts: posts,
	})
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func developmentGetHandler(w http.ResponseWriter, r *http.Request) {
	baseDir, err := os.Getwd()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := template2.Must(template2.ParseFiles(filepath.Join(baseDir, basePathToTemplates, "development.html")))
	err = tmpl.Execute(w, nil)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}