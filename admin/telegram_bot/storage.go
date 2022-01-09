package telegram_bot

import (
	"celme/contacts/back_form"
	"celme/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//Собираем данные полученные ботом
func collectUser(user User, mongoURI, dbName string) error {
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return err
	}
	collection := db.Session.Database(db.Name).Collection("telegram_user")
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"user_id", user.UserId}}
	_, err := collection.UpdateOne(db.Ctx, filter, bson.M{"$set": user}, opts)
	if err != nil {
		return err
	}
	return nil
}

//Ищем пользователя
func findUserUser(userId int, mongoURI, dbName string) (User, error) {
	result := User{}
	crit := bson.M{"user_id": userId}
	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return result, err
	}
	collection := db.Session.Database(db.Name).Collection("telegram_user")
	err := collection.FindOne(db.Ctx, crit).Decode(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	collectionPhones := db.Session.Database(db.Name).Collection("validPhone")
	phone := ValidPhones{}
	err = collectionPhones.FindOne(db.Ctx, bson.M{}).Decode(&phone)
	if err != nil {
		log.Println(err)
		return result, err
	}
	for _, p := range phone.Phones {
		if p == result.PhoneNumber {
			if !result.IsValid {
				result.IsValid = true
				collectUser(result, mongoURI, dbName)
			}
		}
	}
	return result, nil
}

//Ищем пользователлей для оповещения
func findUsersForSubscribe(mongoURI, dbName string) ([]User, error) {
	results := []User{}

	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return results, err
	}
	collectionPhones := db.Session.Database(db.Name).Collection(
		"validPhone")
	phone := ValidPhones{}
	err := collectionPhones.FindOne(db.Ctx, bson.M{}).Decode(&phone)
	if err != nil {
		log.Println(err)
		return results, err
	}
	crit := bson.M{"phone_number": bson.M{"$in": phone.Phones}}

	collection := db.Session.Database(db.Name).Collection("telegram_user")
	cur, err := collection.Find(db.Ctx, crit)
	if err != nil {
		log.Println(err)
		return results, err
	}
	defer cur.Close(db.Ctx)
	for cur.Next(db.Ctx) {
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return results, err
		}
		results = append(results, elem)
	}

	return results, nil
}

//Ищем сообщения для оповещения
func findMessagesSubscribe(mongoURI, dbName string) ([]back_form.Message, error) {
	results := []back_form.Message{}

	db := storage.New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return results, err
	}

	crit := bson.M{"isViewed": false}

	collection := db.Session.Database(db.Name).Collection("messages")
	cur, err := collection.Find(db.Ctx, crit)
	if err != nil {
		log.Println(err)
		return results, err
	}
	defer cur.Close(db.Ctx)
	for cur.Next(db.Ctx) {
		var elem back_form.Message
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
			return results, err
		}
		results = append(results, elem)
	}
	collection.UpdateMany(db.Ctx, crit, bson.M{"$set": bson.M{"isViewed": true}})

	return results, nil
}

type User struct {
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
	ChatID      int64  `bson:"chat_id" json:"chat_id"`
	FirstName   string ` bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	UserId      int    `bson:"user_id" json:"user_id"`
	IsValid     bool
}
type ValidPhones struct {
	Phones []string `bson:"phones" json:"phones"`
}
