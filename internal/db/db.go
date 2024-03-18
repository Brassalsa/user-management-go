package db

import (
	"context"
	"errors"
	"log"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Url         string
	collections []string
	Ctx         context.Context
	client      *mongo.Client
	database    *mongo.Database
}

// connect to db
func (db *Database) Connect(ctx context.Context, dbName string) error {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI(db.Url)
	db.Ctx = ctx
	// Create a new MongoDB client
	client, err := mongo.Connect(db.Ctx, clientOptions)
	if err != nil {
		return err
	}

	db.client = client

	// Ping the MongoDB server to verify the connection
	err = client.Ping(db.Ctx, nil)
	if err != nil {
		return err
	}

	db.database = client.Database(dbName)
	collection := client.Database(dbName).Collection("users")

	iMUsername := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, iMUsername)
	if err != nil {
		log.Fatal(err)
	}
	iMEmail := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, iMEmail)
	if err != nil {
		log.Fatal(err)
	}

	db.collections = []string{"users"}
	return nil
}

// close connection
func (db *Database) Close() error {
	return db.client.Disconnect(db.Ctx)
}

// add data to collection
func (db *Database) InsertIntoCollection(collectionName string, payload interface{}) (*mongo.InsertOneResult, error) {
	if c := helpers.Contains(db.collections, collectionName); !c {
		return nil, errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)
	return coll.InsertOne(db.Ctx, payload)
}

// delete data from collection
func (db *Database) DeleteFromCollection(collectionName string, filter interface{}) error {
	if c := helpers.Contains(db.collections, collectionName); !c {
		return errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)

	_, err := coll.DeleteOne(db.Ctx, filter)
	return err
}

// find data
func (db *Database) FindOne(collectionName string, filter interface{}) (*mongo.SingleResult, error) {

	if c := helpers.Contains(db.collections, collectionName); !c {
		return nil, errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)
	res := coll.FindOne(db.Ctx, filter)
	return res, nil
}

func (db *Database) ConvertStrToId(s string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}

// update data by id
func (db *Database) UpdateById(collectionName string, id primitive.ObjectID, updateParam interface{}) (*mongo.UpdateResult, error) {
	if c := helpers.Contains(db.collections, collectionName); !c {
		return nil, errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)
	res, err := coll.UpdateByID(db.Ctx, id, bson.D{{
		Key:   "$set",
		Value: updateParam,
	}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// get all the data
func (db *Database) GetAllDocuments(collectionName string, results interface{}) error {
	if c := helpers.Contains(db.collections, collectionName); !c {
		return errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)
	cursor, err := coll.Find(db.Ctx, bson.D{})

	if err != nil {
		return err
	}

	defer cursor.Close(db.Ctx)

	if err := cursor.All(db.Ctx, results); err != nil {
		return err
	}

	return nil
}
