package routing

import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/midlleware"
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"fmt"
	"github.com/gorilla/mux"
	template2 "html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Routing interface {
	Setup()
	Run()
}

type APIRouter struct {
	router *mux.Router
}

type WEBRouter struct {
	router *mux.Router
}

func NewAPIRouter() *APIRouter {
	return &APIRouter{
		router: mux.NewRouter(),
	}
}

func NewWEBRouter() *WEBRouter {
	return &WEBRouter{
		router: mux.NewRouter(),
	}
}

var globalTemplate *template2.Template

func (r *APIRouter)Setup() {
	r.router.HandleFunc("/", indexGetHandler).Methods("GET")
}

func (r *WEBRouter)Setup() {
	r.router.Handle("/", midlleware.LoggerMiddleware(http.HandlerFunc(homeGetHandler))).Methods("GET")
	r.router.Handle("/home", midlleware.LoggerMiddleware(http.HandlerFunc(homeGetHandler))).Methods("GET")
	r.router.Handle("/home", midlleware.LoggerMiddleware(http.HandlerFunc(homePostHandler))).Methods("POST")
	r.router.Handle("/sources", midlleware.LoggerMiddleware(http.HandlerFunc(listOfSourcesGetHandler))).Methods("GET")
	r.router.Handle("/development", midlleware.LoggerMiddleware(http.HandlerFunc(developmentGetHandler))).Methods("GET")
	r.router.Handle("/sysInfo", midlleware.LoggerMiddleware(http.HandlerFunc(systemInfoGetHandler))).Methods("GET")

	baseDir, err := os.Getwd()
	if err != nil {
		loger.PrintError(err)
		return
	}
	staticDir := http.Dir(filepath.Join(baseDir, "src", "assets"))
	staticFS := http.FileServer(staticDir)

	r.router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", midlleware.LoggerMiddleware(staticFS)))

	filesDir := http.Dir(filepath.Join(baseDir, "src", "files"))
	filesFS := http.FileServer(filesDir)
	r.router.PathPrefix("/files/").
		Handler(http.StripPrefix("/files/", midlleware.LoggerMiddleware(filesFS)))

	globalTemplate, err = template2.New("sharingToMe").Funcs(template2.FuncMap{
		"now": time.Now,
		"timeZoneOffset": utility.GetTimeZoneOffsetFromGTW,
		"n2": func(number float64) string {
			return fmt.Sprintf("%.2f", number)
		},
		"ng": utility.NumberGroup,
	}).ParseGlob(filepath.Join(baseDir, defaults.BasePathToTemplates))
	if err != nil {
		panic(err)
	}
}

func (r *APIRouter)Run() {
	host, port := configuration.Configuration.ApiServer.Host, configuration.Configuration.ApiServer.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		loger.PrintError(err)
	} else {
		defer listener.Close()
		loger.PrintInfo("Successfully started the API server on the port %s:%d", host, port)
	}
	if err := http.Serve(listener, r.router); err != nil {
		loger.PrintError(err)
	}

}

func (r *WEBRouter)Run() {
	host, port := configuration.Configuration.WebServer.Host, configuration.Configuration.WebServer.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	} else {
		defer listener.Close()
		loger.PrintInfo("Successfully started the WEB server on the port %s:%d", host, port)
	}
	if err := http.Serve(listener, r.router); err != nil {
		panic(err)
	}

}