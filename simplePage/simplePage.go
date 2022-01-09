package simplePage

import (
	"celme/config"
	simplePage "celme/simplePage/simplePage_info"
	"net/http"
	"path"
)

func Register(conf config.Config, mux *http.ServeMux) {
	mux.HandleFunc(
		path.Join(conf.Base, "page", "show"),
		simplePage.RestSimplePageHandler(conf.MongoURI, conf.DBName))

}
