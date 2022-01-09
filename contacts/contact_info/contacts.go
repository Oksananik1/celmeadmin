package contact_info

import (
	helper "celme/helpers"
	"celme/storage"
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

// RestContactsHandler
// возвращает обработчик запросов отображения контактов возвращает JSON
func RestContactsHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db := storage.New(mongoURI, dbName)
		defer db.Close()
		if err := db.Dial(); err != nil {
			helper.HandleError(w, errors.Wrapf(err, "Dial"))
			return
		}
		collection := db.Session.Database(db.Name).Collection("contacts")
		result := bson.M{}
		err := collection.FindOne(db.Ctx, bson.M{}).Decode(&result)
		if err != nil {
			helper.HandleError(w, errors.Wrapf(err, "FindOne & Decode"))
			return
		}

		js, err := json.Marshal(result)
		if err != nil {
			helper.HandleError(w, errors.Wrapf(err, "Marshal"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	}
}
