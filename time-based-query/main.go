package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Entry struct {
	DeviceID  string    `bson:"deviceId"`
	Value     float64   `bson:"value"`
	CreatedAt time.Time `bson:"createdAt"`
}

func main() {
	ctx := context.Background()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("exampleDB").Collection("entries")

	// 🔍 Create index on createdAt for fast time queries
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "createdAt", Value: 1}},
	})
	if err != nil {
		panic(err)
	}

	// ✅ Insert one sample entry
	now := time.Now()
	entry := Entry{DeviceID: "abc-123", Value: 42.5, CreatedAt: now}
	_, _ = collection.InsertOne(ctx, entry)

	// 🔹 Query last 1000 entries
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(1000)
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	var lastEntries []Entry
	_ = cursor.All(ctx, &lastEntries)
	fmt.Printf("Last 1000 entries: %d found\n", len(lastEntries))

	// 🔹 Query entries from past 15 minutes
	fifteenMinsAgo := now.Add(-15 * time.Minute)
	filter := bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: fifteenMinsAgo}}}}

	cursor, err = collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var recentEntries []Entry
	_ = cursor.All(ctx, &recentEntries)
	fmt.Printf("Entries from past 15 minutes: %d found\n", len(recentEntries))
}
