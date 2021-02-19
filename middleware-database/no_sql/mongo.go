package no_sql

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"context"
	"fmt"
)

var (
	mongoDb *mongo.Client
)

func _init(){
	if mongoDb == nil {
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:MongoDB2019!@localhost:27017/?authSource=admin"))
		if err != nil {
			log.Fatal(err)
		}
		mongoDb = client
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		//defer client.Disconnect(ctx)
	}
}

func GetMongoDatabase() *mongo.Client {
	_init()
	return mongoDb
}

func MongoInsert(database string, colName string, entity interface{}) interface{} {
	_init()
	collection := mongoDb.Database(database).Collection(colName)
	insertResult, err := collection.InsertOne(context.TODO(), entity)
	if err != nil{
		fmt.Println("Error to insert data in Mongo")
	}
	return insertResult.InsertedID
}