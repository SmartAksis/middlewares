package broker_consumer

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	dataBaseName = "TRACKING"
	BrokerStatusCollectionName = "BROKER_STATUS"
)

type onSuccess func(string, []byte, string, string)

type CallbackConsumerInterface interface {
	OnSuccess(queueUrl string, data []byte, messageId string, receivedHandle string)
}

func (c *callbackConsumer) OnSuccess(queueUrl string, data []byte, messageId string, receivedHandle string){
	c.success(queueUrl, data, messageId, receivedHandle);
}

type callbackConsumer struct {
	success onSuccess
}

func DoneCallbackConsumer(mongoDb *mongo.Client, application string) *callbackConsumer{
	return &callbackConsumer{
		success: func(queueName string, data []byte, messageId string, receivedHandle string) {
			collection := mongoDb.Database(dataBaseName).Collection(BrokerStatusCollectionName)
			var assync AssyncProcess

			filter := bson.D{
				{"hashcode", messageId},
			}

			result := collection.FindOne(context.TODO(), filter).Decode(&assync)

			if assync.Data != nil && result != nil {
				doneItem:=Done(application)
				assync.Items=append(assync.Items, *doneItem)
				_, err := collection.UpdateOne(context.TODO(), filter, bson.D{
					{"$set", bson.D{{"items", assync.Items}}},
				})
				if err != nil {
					fmt.Println("Erro to add item")
					fmt.Println(err)
				}
			}
		},
	}
}

func StartTransaction(mongoDb *mongo.Client, hashcode string, application string, data interface{}) {
	assync := NewSqsAssyncProcess(hashcode, "accounts", data)
	insert(mongoDb, assync)
}

func DoneTransaction(mongoDb *mongo.Client, hashcode string, application string) {
	collection := mongoDb.Database(dataBaseName).Collection(BrokerStatusCollectionName)
	var assync AssyncProcess

	filter := bson.D{
		{"hashcode", hashcode},
	}

	result := collection.FindOne(context.TODO(), filter).Decode(&assync)

	if assync.Data != nil && result != nil {
		doneItem:=Done(application)
		assync.Items=append(assync.Items, *doneItem)
		_, err := collection.UpdateOne(context.TODO(), filter, bson.D{
			{"$set", bson.D{{"items", assync.Items}}},
		})
		if err != nil {
			fmt.Println("Erro to add item")
			fmt.Println(err)
		}
	}
}

func getBrokerStatusCollection(mongoDb *mongo.Client) *mongo.Collection{
	return mongoDb.Database(dataBaseName).Collection(BrokerStatusCollectionName)
}

func insert(mongoDb *mongo.Client, entity interface{}) interface{} {
	collection := mongoDb.Database(dataBaseName).Collection(BrokerStatusCollectionName)
	insertResult, err := collection.InsertOne(context.TODO(), entity)
	if err != nil{
		fmt.Println("Error to insert data in Mongo")
	}
	return insertResult.InsertedID
}