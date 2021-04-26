package routing

import (
	"IosifSuzuki/sharingToMe/internal/configuration"
	"IosifSuzuki/sharingToMe/internal/midlleware"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"path/filepath"
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

func (r *APIRouter)Setup() {
	r.router.HandleFunc("/", indexGetHandler).Methods("GET")
}

func (r *WEBRouter)Setup() {
	r.router.Handle("/", midlleware.LoggerMiddleware(http.HandlerFunc(homeGetHandler))).Methods("GET")
	r.router.Handle("/home", midlleware.LoggerMiddleware(http.HandlerFunc(homeGetHandler))).Methods("GET")
	r.router.Handle("/home", midlleware.LoggerMiddleware(http.HandlerFunc(homePostHandler))).Methods("POST")
	r.router.Handle("/sources", midlleware.LoggerMiddleware(http.HandlerFunc(listOfSourcesGetHandler))).Methods("GET")
	r.router.Handle("/development", midlleware.LoggerMiddleware(http.HandlerFunc(developmentGetHandler))).Methods("GET")
	baseDir, err := os.Getwd()
	if err != nil {
		loger.PrintError(err)
		return
	}
	dir := http.Dir(filepath.Join(baseDir, "src", "assets"))
	fs := http.FileServer(dir)

	r.router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", midlleware.LoggerMiddleware(fs)))
}

func (r *APIRouter)Run() {
	host, port := configuration.Configuration.ApiServer.Host, configuration.Configuration.ApiServer.Port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		loger.PrintError(err)
	} else {
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
		loger.PrintError(err)
	} else {
		loger.PrintInfo("Successfully started the WEB server on the port %s:%d", host, port)
	}
	if err := http.Serve(listener, r.router); err != nil {
		loger.PrintError(err)
	}

}