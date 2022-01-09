package productsInfo

import (
	helper "celme/helpers"
	"celme/storage"
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

// RestListHandler
// возвращает обработчик запросов отображения товаров возвращает JSON
func RestListHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		grKeys, ok := r.URL.Query()["gr"]
		if !ok || len(grKeys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'gr' is empty"))
			return
		}
		groupFilter := grKeys[0]
		results, err := loadList(mongoURI, dbName, groupFilter)
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

// RestGetHandler
// возвращает обработчик запросов отображения товара возвращает JSON
func RestGetHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["id"]
		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'id' not found"))
			return
		}
		idStr := keys[0]
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'id' is not MongoObjectId"))
			return
		}
		result, err := loadForEdit(mongoURI, dbName, id)
		if err != nil {
			helper.HandleError(w, errors.Wrapf(err, "Marshal"))
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

func loadForEdit(mongoURI, dbName string, id primitive.ObjectID) (ListItem, error) {
	result := ListItem{}
	crit := bson.M{"_id": id}
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return result, err
	}
	collection := db.Session.Database(db.Name).Collection("products")
	err := collection.FindOne(db.Ctx, crit).Decode(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func loadList(mongoURI, dbName, group string) ([]ListItem, error) {
	results := []ListItem{}
	crit := bson.M{"group": group}
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return results, err
	}
	collection := db.Session.Database(db.Name).Collection("products")
	cur, err := collection.Find(db.Ctx, crit)
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
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Descr      string             `bson:"descr" json:"descr"`
	SmallDescr string             `bson:"small_descr" json:"smallDescr"`
	Price      string             `bson:"price" json:"price"`
	Photo      string             `bson:"photo"  json:"photo"`
	Group      string             `bson:"group"  json:"group"`
	Order      string             `bson:"order"  json:"order"`
	Feature    []Features         `bson:"feature" json:"feature"`
}
type Features struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}
