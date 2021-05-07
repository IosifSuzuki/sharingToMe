package routing

import (
	"IosifSuzuki/sharingToMe/internal/JWT"
	"IosifSuzuki/sharingToMe/internal/dbManager"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/fileManager"
	"IosifSuzuki/sharingToMe/internal/gitHubManager"
	"IosifSuzuki/sharingToMe/internal/ipManager"
	"IosifSuzuki/sharingToMe/internal/models"
	"IosifSuzuki/sharingToMe/internal/sysInfoManager"
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"encoding/json"
	"net/http"
	"path/filepath"
)


// API Handlers

func welcomeGetHandler(w http.ResponseWriter, _ *http.Request) {
	var notification = models.Notification{
		Title:       "Welcome to us",
		Description: "This service provide send & store date to server",
		Error: models.Error{
			Message: defaults.RequestCompletedSuccessfully,
			Code:    http.StatusOK,
		},
	}
	w.Header().Set(defaults.ContentType, "application/json")
	err := json.NewEncoder(w).Encode(notification)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		})
		return
	}
}

func signInConsumerPostHandler(w http.ResponseWriter, r *http.Request) {
	var (
		credential models.Credential
		err        = json.NewDecoder(r.Body).Decode(&credential)
	)
	w.Header().Set(defaults.ContentType, "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	consumer, err := dbManager.FetchConsumer(credential)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		})
		return
	}
	token, err := JWT.GenerateToken(*consumer)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(struct {
		Token string       `json:"token"`
		Error models.Error `json:"error"`
	}{
		Token: token,
		Error: models.Error{
			Message: defaults.RequestCompletedSuccessfully,
			Code:    http.StatusBadRequest,
		},
	})
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func signUpConsumerPostHandler(w http.ResponseWriter, r *http.Request) {
	var (
		consumer models.Consumer
		err      = json.NewDecoder(r.Body).Decode(&consumer)
	)
	w.Header().Set(defaults.ContentType, "application/json")
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		})
		return
	}
	consumer.Id = defaults.NewId
	ok, errorMessage := consumer.ValidateForm()
	if !ok {
		loger.PrintError(err)
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: errorMessage,
				Code:    http.StatusForbidden,
			},
		})
		return
	} else {
		consumer.Password, err = utility.CreateHashPassword(consumer.Password)
		if err != nil {
			loger.PrintError(err)
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(struct {
				Error models.Error `json:"error"`
			}{
				Error: models.Error{
					Message: err.Error(),
					Code:    http.StatusForbidden,
				},
			})
		}
	}
	ip, err := utility.GetIP(r)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		})
		return
	}
	isExist, err := dbManager.IsExistConsumer(consumer)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isExist {
		loger.PrintError(err)
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: defaults.ConsumerAlreadyExist,
				Code:    http.StatusConflict,
			},
		})
		return
	}
	ipInfo, err := ipManager.GetIpInfo(ip)
	if ipInfo != nil {
		consumer.IpInfo = *ipInfo
	}
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
		return
	}
	err = dbManager.SaveConsumer(consumer)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(struct {
			Error models.Error `json:"error"`
		}{
			Error: models.Error{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
		return
	}
	err = json.NewEncoder(w).Encode(struct {
		Error models.Error `json:"error"`
	}{
		Error: models.Error{
			Message: defaults.RequestCompletedSuccessfully,
			Code:    http.StatusOK,
		},
	})
}

// WEB

func homeGetHandler(w http.ResponseWriter, r *http.Request) {
	err := globalTemplate.ExecuteTemplate(w, "home.gohtml", struct {
		Title  string
		Errors map[string][]string
	}{
		Title: "Home",
	})
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func homePostHandler(w http.ResponseWriter, r *http.Request) {
	const (
		nickNameKey    = "nickname"
		emailKey       = "email"
		descriptionKey = "description"
		fileKey        = "file"
		maxSizeFile    = 150 * 1024 * 1024
	)
	var err = r.ParseForm()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.ParseMultipartForm(maxSizeFile)
	var (
		nicknameValue    = r.FormValue(nickNameKey)
		emailValue       = r.FormValue(emailKey)
		descriptionValue = r.FormValue(descriptionKey)
	)
	file, fileHeader, _ := r.FormFile(fileKey)
	var newPost = models.NewPost{
		Nickname:      nicknameValue,
		Email:         emailValue,
		Description:   descriptionValue,
		File:          &file,
		FileHeader:    fileHeader,
		MaxSizeOfFile: maxSizeFile,
	}
	ip, err := utility.GetIP(r)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var errors = newPost.ValidatePostForm()
	isAllow, err := dbManager.AllowCreatePost(ip)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAllow {
		errors["AccessRefused"] = make([]string, 0)
		errors["AccessRefused"] = append(errors["AccessRefused"], "File upload limit exceeded")
	}
	if len(errors) > 0 {
		err := globalTemplate.ExecuteTemplate(w, "home.gohtml", struct {
			Title  string
			Errors map[string][]string
		}{
			Title:  "Home",
			Errors: errors,
		})
		if err != nil {
			loger.PrintError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	defer file.Close()
	var fileExtension = filepath.Ext(fileHeader.Filename)
	if len(fileExtension) == 0 {
		fileExtension = ".mp3"
	}
	newFilepath, err := fileManager.SaveMediaFile(file, fileExtension)
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	publisherId, err := utility.GetIP(r)
	if err != nil {
		publisherId = "0.0.0.0"
	}
	var publisher = models.Publisher{
		Id:       defaults.NewId,
		Nickname: nicknameValue,
		Email:    emailValue,
	}
	isExist, err := dbManager.IsExistPublisher(publisher)
	if !isExist {
		ipInfo, err := ipManager.GetIpInfo(publisherId)
		if err != nil {
			loger.PrintError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if ipInfo != nil {
			publisher.IpInfo = *ipInfo
		}
	}
	var post = models.Post{
		Id:          defaults.NewId,
		Description: descriptionValue,
		FilePath:    *newFilepath,
		Publisher:   &publisher,
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
	err = globalTemplate.ExecuteTemplate(w, "listOfSources.gohtml", struct {
		Title string
		Posts []models.Post
	}{
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
	commitInfos, err := gitHubManager.FetchCommitInfos()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = globalTemplate.ExecuteTemplate(w, "development.gohtml", struct {
		Title       string
		CommitInfos []models.CommitInfo
	}{
		Title:       "Development",
		CommitInfos: commitInfos,
	})
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func systemInfoGetHandler(w http.ResponseWriter, r *http.Request) {
	var sysInfo, err = sysInfoManager.GetSysInfo()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = globalTemplate.ExecuteTemplate(w, "sysInfo.gohtml", struct {
		Title   string
		SysInfo models.SysInfo
	}{
		Title:   "System Info",
		SysInfo: *sysInfo,
	})
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func clearDataGetHandler(w http.ResponseWriter, r *http.Request) {
	filePaths, err := dbManager.ClearOldData()
	if err != nil {
		loger.PrintError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, filePath := range filePaths {
		err = fileManager.RemoveFile(filePath)
		if err != nil {
			loger.PrintError(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/sources", http.StatusSeeOther)
}
