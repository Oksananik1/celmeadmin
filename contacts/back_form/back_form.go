package back_form

import (
	"celme/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

//Собираем данные полученные c сайта
func SaveHandler(mongoURI, dbName string) func(http.ResponseWriter,
	*http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// function body of a http.HandlerFunc
		message := Message{}
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Error on parse form"))
			return
		}

		message.ID = primitive.NewObjectID()

		names := r.PostFormValue("name")
		if len(names) < 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Url Param 'name' is empty"))
			return
		}
		message.Name = names
		message.Phone = r.PostFormValue("phone")
		message.Email = r.PostFormValue("email")
		message.Message = r.PostFormValue("message")
		message.IsViewed = false

		collectMessage(message, mongoURI, dbName)
		http.Redirect(w, r, "/feedback", 302)
	}
}
func collectMessage(message Message, mongoURI, dbName string) error {
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return err
	}
	collection := db.Session.Database(db.Name).Collection("messages")
	_, err := collection.InsertOne(db.Ctx, message)
	if err != nil {
		return err
	}
	return nil
}

// Message описывает элемент списка спортсменов, загружаемый из БД
type Message struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string             `bson:"name" json:"name"`
	Phone    string             `bson:"phone" json:"phone"`
	Email    string             `bson:"email" json:"email"`
	Message  string             `bson:"message" json:"message"`
	IsViewed bool               `bson:"isViewed" json:"isViewed"`
}
