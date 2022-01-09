package contacts

import (
	"celme/config"
	"celme/contacts/back_form"
	"celme/contacts/contact_info"
	"net/http"
	"path"
)

func Register(conf config.Config, mux *http.ServeMux) {
	mux.HandleFunc(
		path.Join(conf.Base, "contacts", "show"),
		contact_info.RestContactsHandler(conf.MongoURI, conf.DBName))
	mux.HandleFunc(
		path.Join(conf.Base, "sendmail"),
		back_form.SaveHandler(conf.MongoURI, conf.DBName))

}
