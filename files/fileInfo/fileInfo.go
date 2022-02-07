package fileInfo

import (
	helper "celme/helpers"
	"celme/storage"
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

// RestListHandler
// возвращает обработчик запросов отображения товаров возвращает JSON
func RestListHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		results, err := loadList(mongoURI, dbName)
		if err != nil {
			helper.HandleError(w, errors.Wrapf(err, "loadList"))
			return
		}
		js, err := json.Marshal(results)
		if err != nil {
			helper.HandleError(w, errors.Wrapf(err, "Marshal"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	}
}

func loadList(mongoURI, dbName string) ([]ListItem, error) {
	results := []ListItem{}
	crit := bson.M{}
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return results, err
	}
	collection := db.Session.Database(db.Name).Collection("files")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"order", 1}})
	cur, err := collection.Find(db.Ctx, crit, findOptions)

	if err != nil {
		log.Println(err)
		return results, err
	}
	defer cur.Close(db.Ctx)
	for cur.Next(db.Ctx) {
		var elem ListItem
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return results, err
		}
		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return results, nil
}

// ListItem описывает элемент списка спортсменов, загружаемый из БД
type ListItem struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string             `bson:"name" json:"name"`
	FilePath string             `bson:"filePath" json:"filePath"`
	Descr    string             `bson:"descr" json:"descr"`
	Order    string             `bson:"order"  json:"order"`
}
