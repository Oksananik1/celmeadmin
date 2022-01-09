package storage

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoDB описывает репозиторий доступа к данным
type MongoDB struct {
	mongoURI string
	Name     string
	Session  *mongo.Client
	Ctx      context.Context
}

// New создаёт экземпляр репозитория для доступа к данным
func New(mongoURI, dbName string) *MongoDB {
	return &MongoDB{mongoURI: mongoURI, Name: dbName}
}

// Close закрывает соединение с базой данных
func (db *MongoDB) Close() {
	if db.Session != nil {
		db.Session.Disconnect(db.Ctx)
		db.Session = nil
	}
}

func (db *MongoDB) Dial() error {
	// Set client options
	clientOptions := options.Client().ApplyURI(db.mongoURI)

	// Connect to MongoDB
	db.Ctx = context.TODO()
	client, err := mongo.Connect(db.Ctx, clientOptions)

	if err != nil {
		return errors.Wrap(err, "Not connected to mongodb")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return errors.Wrap(err, "Not ping to mongodb")
	}
	db.Session = client
	return nil
}

type SimpleQueryResult struct {
	Meta struct {
		Limit  int64 `json:"limit"`
		Offset int64 `json:"offset"`
		Total  int64 `json:"total_count"`
	} `json:"meta"`
	Objects []bson.M `json:"objects"`
}
type Query struct {
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
	Projection bson.M
	Sort       bson.M
	Query      bson.M `json:"query"`
}

func DefaultQuery() Query {
	query := Query{}
	query.Limit = 20
	query.Offset = 0
	query.Query = bson.M{}
	query.Projection = bson.M{}
	return query
}
func SimpleQuery(mongoURI, dbName, collectionName string,
	query Query) (SimpleQueryResult, error) {
	results := SimpleQueryResult{}
	result := []bson.M{}
	db := New(mongoURI, dbName)
	defer db.Close()
	if err := db.Dial(); err != nil {
		return results, err
	}
	collection := db.Session.Database(db.Name).Collection(collectionName)
	optionFind := options.FindOptions{
		Limit:      &query.Limit,
		Skip:       &query.Offset,
		Projection: &query.Projection,
		Sort:       &query.Sort,
	}
	cur, err := collection.Find(db.Ctx, query.Query, &optionFind)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "Aggregate error on sportsman."+
			"storage.Fetch "))
		log.Println(err)
		return results, err
	}
	countTotal, err := collection.CountDocuments(db.Ctx, query.Query)
	if err != nil {
		fmt.Println(errors.Wrapf(err, "Aggregate error on sportsman."+
			"storage.Fetch "))
		log.Println(err)
		return results, err
	}
	defer cur.Close(db.Ctx)
	for cur.Next(db.Ctx) {
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(errors.Wrapf(err, "Aggregate error on sportsman."+
				"storage.Fetch "))
			log.Println(err)
			return results, err
		}
		result = append(result, elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(errors.Wrapf(err, "Aggregate error on sportsman."+
			"storage.Fetch "))
		log.Println(err)
	}
	results.Meta.Offset = query.Offset
	results.Meta.Limit = query.Limit
	results.Meta.Total = countTotal
	results.Objects = result
	return results, nil

}
