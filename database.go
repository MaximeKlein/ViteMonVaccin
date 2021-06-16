package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDb struct {
	client mongo.Client
	ctx    context.Context
}

func (db mongoDb) init() {
	/*
	   Connect to my cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	db.client = *client
	if err != nil {
		log.Fatal(err)
	}
	db.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = db.client.Connect(db.ctx)
	if err != nil {
		log.Fatal(err)
	}

	/*
	   List databases
	*/
	databases, err := db.client.ListDatabaseNames(db.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func (db mongoDb) close() {
	db.client.Disconnect(db.ctx)
}
