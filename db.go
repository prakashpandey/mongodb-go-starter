package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var Database *mongo.Database

// IndexType is the type of index ASC or DESC
type IndexType int32

const (
	// SessionCollection in which sessions are stored
	pcpCollection = "pcp"
)

const (
	asc  IndexType = 1
	desc IndexType = -1
)

// createUniqueIndexOnSessionsCollection will panic if creating an index on session collection fail
func createUniqueIndexOnPCPCollection() {
	index := createIndexModel("name", "name", asc, true, 0)
	mustCreateIndex(index, Database.Collection(pcpCollection))
}

// createIndexModel creates mongo.IndexModel object
func createIndexModel(key string, indexName string, indexType IndexType, unique bool, ttl int32) mongo.IndexModel {
	keys := bsonx.Doc{
		{
			Key:   key,
			Value: bsonx.Int32(int32(indexType)),
		},
	}
	index := mongo.IndexModel{}
	index.Keys = keys
	index.Options = &options.IndexOptions{
		Name:               &indexName,
		Unique:             &unique,
		ExpireAfterSeconds: &ttl,
	}
	return index
}

// mustCreateIndex will panic if creating an index on given collection fail
func mustCreateIndex(index mongo.IndexModel, c *mongo.Collection) {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	if _, err := c.Indexes().CreateOne(context.Background(), index, opts); err != nil {
		panic(fmt.Sprintf("error while applying index to collection[%s], error[%s]", c.Name(), err.Error()))
	}
}

func connectMongoDB(connectionString string) (*mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	return mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
}

// FindUserByUsername returns user, error if user do not exist
func FindUserByUsername(name string) (u *User, err error) {
	collection := Database.Collection(pcpCollection)
	if err = collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&u); err != nil {
		return nil, err
	}
	return u, nil
}

func update(u *User) {
	log.Printf("Updating data")
	filter := bson.M{
		"name": u.Name,
	}
	update := bson.M{
		"$set": bson.M{
			"age":   u.Age,
			"phone": u.Phone,
			"city":  u.City,
		},
	}
	upsert := false
	opts := &options.UpdateOptions{Upsert: &upsert}
	collection := Database.Collection(pcpCollection)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data updated")
}

func insert(u *User) {
	log.Printf("Inserting data")
	collection := Database.Collection(pcpCollection)
	_, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data inserted")
}
