package db

import (
	"context"
	"errors"

	"github.com/Brassalsa/user-management-go/pkg/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Url         string
	Collections []string
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

	for _, val := range db.Collections {
		err := db.database.CreateCollection(db.Ctx, val)
		if err != nil {
			return err
		}

	}

	return nil
}

// close connection
func (db *Database) Close() error {
	return db.client.Disconnect(db.Ctx)
}

// add data to collection
func (db *Database) InsertIntoCollection(collectionName string, payload interface{}) (*mongo.InsertOneResult, error) {
	c := helpers.Contains(db.Collections, collectionName)
	if !c {
		return nil, errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)
	return coll.InsertOne(db.Ctx, payload)
}

// delete data from collection
func (db *Database) DeleteFromCollection(collectionName string, filter interface{}) error {
	c := helpers.Contains(db.Collections, collectionName)
	if !c {
		return errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)

	_, err := coll.DeleteOne(db.Ctx, filter)
	return err
}

// find data
func (db *Database) FindOne(collectionName string, filter interface{}) (*mongo.SingleResult, error) {
	c := helpers.Contains(db.Collections, collectionName)
	if !c {
		return nil, errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(collectionName)
	res := coll.FindOne(db.Ctx, filter)
	return res, nil
}
