package videos

import (
	"celme/blank"
	helper "celme/helpers"
	"celme/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/url"
	"path"
)

// ListHandler
// возвращает обработчик запросов отображения списка товаров
func ListHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	base := "/blank/views/videos/list"
	templates := []string{
		path.Join("/blank/views", "pure"),
		"content", "list",
	}
	return func(w http.ResponseWriter, r *http.Request) {
		getTmpl := blank.TemplateLoader(base, helper.Assets)
		tmpl := getTmpl(templates...)
		data := newRenderList(r)
		data.Items, _ = loadList(mongoURI, dbName)
		tmpl.ExecuteTemplate(w, "pure", data)
	}
}

func newRenderList(r *http.Request) *renderList {
	blankData := blank.NewRenderData()
	blankData.Title = "Список видео"
	blankData.Container = "container"

	return &renderList{
		RenderData: blankData,
	}
}

type renderList struct {
	blank.RenderData
	Items  []ListItem
	Active url.Values
}

// ListItem описывает элемент списка продуктов, загружаемый из БД
type ListItem struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Name  string             `bson:"name" json:"name"`
	Descr string             `bson:"descr" json:"descr"`
	Order string             `bson:"order"  json:"order"`
}

func loadList(mongoURI, dbName string) ([]ListItem, error) {
	results := []ListItem{}
	crit := bson.M{}
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return results, err
	}
	collection := db.Session.Database(db.Name).Collection("videos")
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
