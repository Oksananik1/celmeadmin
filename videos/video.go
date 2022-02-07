package videos

import (
	"celme/config"
	"celme/videos/videoInfo"
	"net/http"
	"path"
)

func Register(conf config.Config, mux *http.ServeMux) {
	mux.HandleFunc(
		path.Join(conf.Base, "video", "list"),
		videoInfo.RestListHandler(conf.MongoURI, conf.DBName))

}
