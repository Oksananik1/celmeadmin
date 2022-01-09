package simplePage

import (
	helper "celme/helpers"
	"celme/storage"
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

// RestShipingHandler
// возвращает обработчик запросов отображения контактов возвращает JSON
func RestSimplePageHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'key' is missing"))
			return
		}
		key := keys[0]

		db := storage.New(mongoURI, dbName)
		defer db.Close()
		if err := db.Dial(); err != nil {
			helper.HandleError(w, errors.Wrapf(err, "Dial"))
			return
		}
		collection := db.Session.Database(db.Name).Collection("shipping")
		result := bson.M{}
		err := collection.FindOne(db.Ctx, bson.M{"_id": key}).Decode(&result)
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
