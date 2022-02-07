package files

import (
	"celme/blank"
	helper "celme/helpers"
	"celme/storage"
	"errors"
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
	base := "/blank/views/files/form"
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
			data.Item = ListItem{
				ID:    primitive.NewObjectID(),
				Order: "10",
			}
		}

		err := tmpl.ExecuteTemplate(w, "pure", data)
		if err != nil {
			println(err)
		}
	}
}

func newRenderEdit(r *http.Request) *renderEditFile {

	blankData := blank.NewRenderData()
	blankData.Title = "Добавление/Редактирование файла"
	blankData.Container = "container"

	return &renderEditFile{
		RenderData: blankData,
	}
}

type renderEditFile struct {
	blank.RenderData
	Item ListItem
}

func loadForEdit(mongoURI, dbName string, id primitive.ObjectID) (ListItem, error) {
	result := ListItem{}
	crit := bson.M{"_id": id}
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return result, err
	}
	collection := db.Session.Database(db.Name).Collection("files")
	err := collection.FindOne(db.Ctx, crit).Decode(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

// SaveHandler
// возвращает обработчик запросов сохраниения товара
func SaveHandler(mongoURI, dbName, filePath string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// function body of a http.HandlerFunc
		fileInf := ListItem{}
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
		fileInf.ID = id

		orders := r.MultipartForm.Value["order"]
		if len(orders) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'name' is empty"))
			return
		}

		fileInf.Order = orders[0]

		descrs := r.MultipartForm.Value["descr"]
		if len(descrs) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'descr' is empty"))
			return
		}

		fileInf.Descr = descrs[0]

		file, handler, err := r.FormFile(
			"file") //retrieve the file from form data
		//replace file with the key your sent your image with
		if err == nil {

			defer file.Close() //close the file when we finish
			//this is path which  we want to store the file
			dirPath := "celmeapi/storage/files/" + fileInf.ID.Hex() + "/"
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
			fileInf.FilePath = "/" + dirPath + handler.Filename
			fileInf.Name = handler.Filename
		} else {
			fileNames := r.MultipartForm.Value["fileName"]
			name := r.MultipartForm.Value["name"]
			if len(fileNames) > 0 {
				fileInf.FilePath = fileNames[0]
				fileInf.Name = name[0]
			}

		}
		savetoDB(mongoURI, dbName, fileInf)
		http.Redirect(w, r, "/celmeadmin/files/list", 302)
	}
}

func savetoDB(mongoURI, dbName string, item ListItem) error {
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return err
	}
	collection := db.Session.Database(db.Name).Collection("files")
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
		http.Redirect(w, r, "/celmeadmin/files/list", 302)
	}
}

func deleteFromDB(mongoURI, dbName string, id primitive.ObjectID) error {
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return err
	}
	collection := db.Session.Database(db.Name).Collection("files")
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
