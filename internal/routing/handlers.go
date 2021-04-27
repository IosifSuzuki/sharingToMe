package routing

import (
	"IosifSuzuki/sharingToMe/internal/dbManager"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/fileManager"
	"IosifSuzuki/sharingToMe/internal/models"
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
)

const (
	contentType = "Content-Type"
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

	err := globalTemplate.ExecuteTemplate(w, "home.gohtml", nil)
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
	defer file.Close()
	if len(nicknameValue) == 0 || fileHandler == nil || err != nil {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	var fileExtension = filepath.Ext(fileHandler.Filename)
	if len(fileExtension) == 0 {
		fileExtension = ".mp3"
	}
	newFilepath, err := fileManager.SaveMediaFile(file, fileExtension)
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
		Ip: publisherId,
		Flag: flagURL,
		Latitude: 48.62828063964844,
		Longitude: 22.514659881591797,
	}
	var post = models.Post{
		Id: defaults.NewId,
		Description: descriptionValue,
		FilePath: *newFilepath,
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
	posts, err := dbManager.ReadPosts()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = globalTemplate.ExecuteTemplate(w, "listOfSources.gohtml", struct{
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
	var err = globalTemplate.ExecuteTemplate(w, "development.gohtml", nil)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}