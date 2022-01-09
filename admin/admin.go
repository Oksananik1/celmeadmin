package admin

import (
	"celme/admin/auth"
	"celme/admin/products"
	"celme/config"
	"net/http"
	"path"
)

func Register(conf config.Config, mux *http.ServeMux) {
	mux.HandleFunc(
		path.Join("/celmeadmin", "login"),
		auth.LoginHandler(conf.MongoURI, conf.DBName))
	mux.HandleFunc(
		path.Join("/celmeadmin", "product", "list"),
		products.ListHandler(conf.MongoURI, conf.DBName))
	mux.HandleFunc(
		path.Join("/celmeadmin"),
		products.ListHandler(conf.MongoURI, conf.DBName))

	mux.HandleFunc("/celmeadmin/",
		products.ListHandler(conf.MongoURI, conf.DBName))

	mux.HandleFunc(
		path.Join("/celmeadmin", "product", "edit"),
		products.ShowHandler(conf.MongoURI, conf.DBName))
	mux.HandleFunc(
		path.Join("/celmeadmin", "product", "save"),
		products.SaveHandler(conf.MongoURI, conf.DBName, conf.FilePath))
	mux.HandleFunc(
		path.Join("/celmeadmin", "product", "delete"),
		products.DeleteHandler(conf.MongoURI, conf.DBName))

}
