package web

import (
	"embed"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var dist embed.FS

type UI struct {
	router *mux.Router
}

func NewUI() *UI {
	router := mux.NewRouter().StrictSlash(true)
	var assets, _ = fs.Sub(dist, "dist")
	if err := fs.WalkDir(
		dist, "dist", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				router.Methods("GET").Path("/ui/" + strings.TrimPrefix(path, "dist/")).HandlerFunc(
					func(w http.ResponseWriter, req *http.Request) {
						http.ServeFileFS(w, req, assets, strings.TrimPrefix(path, "dist"))
					},
				)
			}
			return nil
		},
	); err != nil {
		panic(err)
	}
	router.Methods("GET").PathPrefix("/ui/").HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			http.ServeFileFS(w, req, assets, "/index.html")
		},
	)
	router.Use(handlers.CompressHandler)
	return &UI{router}
}

func (u *UI) GetRouter() *mux.Router {
	return u.router
}

//
//func (u *UI) Routes() openapi.Routes {
//	var assets, _ = fs.Sub(dist, "dist")
//	return openapi.Routes{
//		"index": {
//			Method:  "GET",
//			Pattern: "/",
//			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
//				w.Header().Set("Location", "/ui/")
//				w.WriteHeader(http.StatusMovedPermanently)
//			},
//		},
//		"ui": {
//			Method:  "GET",
//			Pattern: "/ui/",
//			HandlerFunc: func(w http.ResponseWriter, req *http.Request) {
//				http.StripPrefix("/ui/", http.FileServer(http.FS(assets))).ServeHTTP(w, req)
//			},
//		},
//	}
//}

//var _ openapi.Router = &UI{}
