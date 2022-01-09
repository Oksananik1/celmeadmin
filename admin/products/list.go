package products

import (
	"celme/blank"
	helper "celme/helpers"
	"celme/qst"
	"celme/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/url"
	"path"
)

func getGroups() []string {
	listGroups := []string{
		"Слайсери",
		"М'ясорубки",
		"Куттери",
		"Овочерізки",
		"Обладнання для піцерії",
	}
	return listGroups
}

type Group struct {
	Name   string
	Active string
}

// ListHandler
// возвращает обработчик запросов отображения списка товаров
func ListHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	base := "/blank/views/products/list"
	templates := []string{
		path.Join("/blank/views", "pure"),
		"content", "list",
	}
	return func(w http.ResponseWriter, r *http.Request) {
		getTmpl := blank.TemplateLoader(base, helper.Assets)
		tmpl := getTmpl(templates...)
		data := newRenderList(r)
		data.Items, _ = loadList(mongoURI, dbName, data.Active.Get("n"))
		tmpl.ExecuteTemplate(w, "pure", data)
	}
}

func newRenderList(r *http.Request) *renderList {
	state := ListState{}
	qst.State(&state, r.URL.Query())
	blankData := blank.NewRenderData()
	blankData.Title = "Список продукции"
	blankData.Container = "container"
	activeFilters := state.Active()
	if !activeFilters.Has("n") {
		activeFilters.Add("n", getGroups()[0])
		state.NameGroup = getGroups()[0]
	}
	groups := []Group{}
	for _, g := range getGroups() {
		gr := Group{
			Name:   g,
			Active: "",
		}
		if gr.Name == activeFilters.Get("n") {
			gr.Active = "active"
		}
		groups = append(groups, gr)
	}

	return &renderList{
		RenderData: blankData,
		State:      state,
		Active:     activeFilters,
		Groups:     groups,
	}
}

type renderList struct {
	blank.RenderData
	Items  []ListItem
	Active url.Values
	State  ListState
	Groups []Group
}

// ListItem описывает элемент списка продуктов, загружаемый из БД
type ListItem struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Descr      string             `bson:"descr" json:"descr"`
	SmallDescr string             `bson:"small_descr" json:"small_descr"`
	Price      string             `bson:"price" json:"price"`
	Photo      string             `bson:"photo"  json:"photo"`
	Group      string             `bson:"group"  json:"group"`
	Order      string             `bson:"order"  json:"order"`
	Feature    []Features         `bson:"feature" json:"feature"`
	FeatureStr []Feature          `json:"-"`
}
type Features struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

// ListState описывает запрос к БД
type ListState struct {
	NameGroup string // n

}

// Encode кодирует состояние в URL query
func (ls ListState) Encode(prefix string, query url.Values) {
	query.Set(prefix+"n", ls.NameGroup)

}

// Decode декодирует состояние из URL query
func (ls *ListState) Decode(prefix string, query url.Values) {
	ls.NameGroup = query.Get(prefix + "n")

}

// Active трансформирует состояние в список активных фильтров
func (ls *ListState) Active() url.Values {
	res := url.Values{}
	if ls.NameGroup != "" {
		res.Add("n", ls.NameGroup)
	}
	return res
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
