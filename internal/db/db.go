package db

import (
	"context"
	"errors"

	"log"

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
func (db *Database) Connect(ctx context.Context, dbName string) {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI(db.Url)
	db.Ctx = ctx
	// Create a new MongoDB client
	client, err := mongo.Connect(db.Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db.client = client

	// Ping the MongoDB server to verify the connection
	err = client.Ping(db.Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.database = client.Database(dbName)

	for _, val := range db.Collections {
		err := db.database.CreateCollection(db.Ctx, val)
		if err != nil {
			log.Fatal(err)
		}

	}
}

// close connection
func (db *Database) Close() error {
	return db.client.Disconnect(db.Ctx)
}

func (db *Database) InsertIntoCollection(name string, payload interface{}) (*mongo.InsertOneResult, error) {
	c := helpers.Contains(db.Collections, name)
	if !c {
		return nil, errors.New("collection doesn't exists")
	}

	coll := db.database.Collection(name)
	return coll.InsertOne(db.Ctx, payload)
}
