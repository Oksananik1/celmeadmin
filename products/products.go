package products

import (
	"celme/config"
	"celme/products/productsInfo"
	"net/http"
	"path"
)

func Register(conf config.Config, mux *http.ServeMux) {
	mux.HandleFunc(
		path.Join(conf.Base, "product", "list"),
		productsInfo.RestListHandler(conf.MongoURI, conf.DBName))
	mux.HandleFunc(
		path.Join(conf.Base, "product", "item"),
		productsInfo.RestGetHandler(conf.MongoURI, conf.DBName))

}
