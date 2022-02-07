package files

import (
	"celme/config"
	"celme/files/fileInfo"
	"net/http"
	"path"
)

func Register(conf config.Config, mux *http.ServeMux) {
	mux.HandleFunc(
		path.Join(conf.Base, "files", "list"),
		fileInfo.RestListHandler(conf.MongoURI, conf.DBName))

}
