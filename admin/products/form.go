package products

import (
	"celme/blank"
	helper "celme/helpers"
	"celme/qst"
	"celme/storage"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// ShowHandler
// возвращает обработчик запросов отображения товара
func ShowHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	base := "/blank/views/products/form"
	templates := []string{
		path.Join("/blank/views", "pure"),
		"content", "form",
	}
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["id"]
		forCreate := false
		if !ok || len(keys[0]) < 1 {
			forCreate = true
		}
		getTmpl := blank.TemplateLoader(base, helper.Assets)
		tmpl := getTmpl(templates...)
		data := newRenderEdit(r)
		if !forCreate {
			idStr := keys[0]
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Url Param 'id' is not MongoObjectId"))
				return
			}
			data.Item, _ = loadForEdit(mongoURI, dbName, id)
		} else {
			grKeys, ok := r.URL.Query()["gr"]
			if !ok || len(grKeys[0]) < 1 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Url Param 'gr' is empty"))
				return
			}
			fEmpty := []Feature{}
			fEmpty = append(fEmpty, Feature{"", ""})
			data.Item = ListItem{
				ID:         primitive.NewObjectID(),
				Group:      grKeys[0],
				FeatureStr: fEmpty,
				Order:      "10",
			}
		}
		groups := []Group{}
		for _, g := range getGroups() {
			gr := Group{
				Name:   g,
				Active: "",
			}
			if gr.Name == data.Item.Group {
				gr.Active = "selected"
			}
			groups = append(groups, gr)
		}
		data.Groups = groups

		err := tmpl.ExecuteTemplate(w, "pure", data)
		if err != nil {
			println(err)
		}
	}
}

func newRenderEdit(r *http.Request) *renderEditProduct {
	state := ListState{}
	qst.State(&state, r.URL.Query())
	blankData := blank.NewRenderData()
	blankData.Title = "Создание/Редактирование товара"
	blankData.Container = "container"

	return &renderEditProduct{
		RenderData: blankData,
	}
}

type renderEditProduct struct {
	blank.RenderData
	Item   ListItem
	Groups []Group
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
	features := []Feature{}
	for _, f := range result.Feature {
		features = append(features, Feature{f.Name, f.Value})
	}

	result.FeatureStr = features
	return result, nil
}

// SaveHandler
// возвращает обработчик запросов сохраниения товара
func SaveHandler(mongoURI, dbName, filePath string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// function body of a http.HandlerFunc
		product := ListItem{}
		r.ParseMultipartForm(32 << 20)

		idStrs := r.MultipartForm.Value["id"]
		if len(idStrs) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'id' is empty"))
			return
		}
		id, err := primitive.ObjectIDFromHex(idStrs[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'id' is not MongoObjectId"))
			return
		}
		product.ID = id

		names := r.MultipartForm.Value["name"]
		if len(names) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'name' is empty"))
			return
		}
		product.Name = names[0]

		orders := r.MultipartForm.Value["order"]
		if len(orders) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'name' is empty"))
			return
		}

		product.Order = orders[0]

		groups := r.MultipartForm.Value["group"]
		if len(names) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'group' is empty"))
			return
		}

		product.Group = groups[0]

		prices := r.MultipartForm.Value["price"]
		if len(names) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'group' is empty"))
			return
		}

		product.Price = prices[0]

		smallDescrs := r.MultipartForm.Value["smallDescr"]
		if len(smallDescrs) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'smallDescr' is empty"))
			return
		}

		product.SmallDescr = smallDescrs[0]

		descrs := r.MultipartForm.Value["descr"]
		if len(smallDescrs) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'descr' is empty"))
			return
		}

		product.Descr = descrs[0]

		features := r.MultipartForm.Value["feature"]
		if len(smallDescrs) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'feature' is empty"))
			return
		}

		fmt.Println(features)
		var arr []Feature
		_ = json.Unmarshal([]byte(features[0]), &arr)
		product.Feature = decodeFeatures(arr)

		file, handler, err := r.FormFile("photoFile") //retrieve the file from form data
		//replace file with the key your sent your image with
		if err == nil {

			defer file.Close() //close the file when we finish
			//this is path which  we want to store the file
			dirPath := "celmeapi/storage/product/" + product.ID.Hex() + "/"
			if err := ensureDir(filePath + dirPath); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("Directory creation failed with error: " + err.Error()))
				return
			}
			f, err := os.OpenFile(filePath+dirPath+handler.Filename,
				os.O_WRONLY|os.O_CREATE, 0777)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			defer f.Close()
			io.Copy(f, file)
			product.Photo = "/" + dirPath + handler.Filename
		} else {
			fileNames := r.MultipartForm.Value["fileName"]
			if len(fileNames) > 0 {
				product.Photo = fileNames[0]
			}

		}
		savetoDB(mongoURI, dbName, product)
		http.Redirect(w, r, "/celmeadmin/product/list?n="+product.Group, 302)
	}
}

type Feature []string

func decodeFeatures(features []Feature) []Features {
	result := []Features{}
	for _, f := range features {
		result = append(result, Features{
			Name:  f[0],
			Value: f[1],
		})
	}
	return result
}

func savetoDB(mongoURI, dbName string, item ListItem) error {
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return err
	}
	collection := db.Session.Database(db.Name).Collection("products")
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", item.ID}}
	_, err := collection.UpdateOne(db.Ctx, filter, bson.M{"$set": item}, opts)
	if err != nil {
		return err
	}
	return nil
}

// DeleteHandler
// возвращает обработчик запросов удаления товара
func DeleteHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["id"]
		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'id' is not found"))
			return
		}
		idStr := keys[0]
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'id' is not MongoObjectId"))
			return
		}
		err = deleteFromDB(mongoURI, dbName, id)
		if err != nil {
			println(err)
		}
		http.Redirect(w, r, "/celmeadmin/product/list", 302)
	}
}

func deleteFromDB(mongoURI, dbName string, id primitive.ObjectID) error {
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return err
	}
	collection := db.Session.Database(db.Name).Collection("products")
	_, err := collection.DeleteOne(db.Ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, 0777)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
