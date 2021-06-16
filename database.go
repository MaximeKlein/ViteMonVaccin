package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDb struct {
	client *mongo.Client
	ctx    context.Context
}

type VaccinationCenter struct {
	Name    string `bson:"name,omitempty"`
	Address string `bson:"address,omitempty"`
	Id      string `bson:"_id,omitempty"`
}

type TimeSlot struct {
	Date     string `bson:"date,omitempty"`
	Time     string `bson:"time,omitempty"`
	BookedBy string `bson:"bookedby,omitempty"`
	Id       string `bson:"_id,omitempty"`
}

func initDb() (*mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/ViteMonVaccin"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx
}

func (db mongoDb) getAllVaccinationCenters() []VaccinationCenter {
	collection := db.client.Database("ViteMonVaccin").Collection("VaccinationCenters")

	cur, err := collection.Find(db.ctx, bson.D{})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer cur.Close(db.ctx)

	var centers []VaccinationCenter
	if err = cur.All(db.ctx, &centers); err != nil {
		panic(err)
	}
	fmt.Println(centers)

	return centers
}

func (db mongoDb) getTimeslots(locationId string, restricted bool) []TimeSlot {
	collection := db.client.Database("ViteMonVaccin").Collection("Timeslots")

	selector := bson.M{"locationId": locationId}
	if restricted {
		selector = bson.M{"locationId": locationId, "bookedby": bson.M{"$exists": false}}
	}
	cur, err := collection.Find(db.ctx, selector)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer cur.Close(db.ctx)

	var centers []TimeSlot
	if err = cur.All(db.ctx, &centers); err != nil {
		panic(err)
	}
	fmt.Println(centers)

	return centers
}

func (db mongoDb) reserveTimeslot(timeSlotId string, bookerName string) error {
	collection := db.client.Database("ViteMonVaccin").Collection("Timeslots")

	id, _ := primitive.ObjectIDFromHex(timeSlotId)

	var timeslot TimeSlot

	err := collection.FindOne(db.ctx, bson.M{"_id": id}).Decode(&timeslot)

	if err != nil || timeslot.Date == "" || timeslot.BookedBy != "" {
		return errors.New("timeslot unavailable")
	}
	result, err := collection.UpdateOne(
		db.ctx,
		bson.M{"_id": id},
		bson.D{{"$set", bson.D{{"bookedby", bookerName}}}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	return nil
}

func (db mongoDb) close() {
	db.client.Disconnect(db.ctx)
}
