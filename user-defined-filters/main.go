package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	DeviceID  string    `bson:"deviceId"`
	A         float64   `bson:"A"`
	B         float64   `bson:"B"`
	CreatedAt time.Time `bson:"createdAt"`
}

// Simple safe parser for filters like "A > 10 AND B < 5"
func parseFilter(filterStr string) bson.D {
	filter := bson.D{}
	conditions := strings.Split(filterStr, "AND")
	for _, cond := range conditions {
		cond = strings.TrimSpace(cond)
		var field string
		var op string
		var value float64
		if strings.Contains(cond, ">=") {
			fmt.Sscanf(cond, "%s >= %f", &field, &value)
			filter = append(filter, bson.E{Key: field, Value: bson.D{{Key: "$gte", Value: value}}})
		} else if strings.Contains(cond, "<=") {
			fmt.Sscanf(cond, "%s <= %f", &field, &value)
			filter = append(filter, bson.E{Key: field, Value: bson.D{{Key: "$lte", Value: value}}})
		} else if strings.Contains(cond, ">") {
			fmt.Sscanf(cond, "%s > %f", &field, &value)
			filter = append(filter, bson.E{Key: field, Value: bson.D{{Key: "$gt", Value: value}}})
		} else if strings.Contains(cond, "<") {
			fmt.Sscanf(cond, "%s < %f", &field, &value)
			filter = append(filter, bson.E{Key: field, Value: bson.D{{Key: "$lt", Value: value}}})
		} else if strings.Contains(cond, "=") {
			fmt.Sscanf(cond, "%s = %f", &field, &value)
			filter = append(filter, bson.E{Key: field, Value: value})
		}
	}
	return filter
}

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("exampleDB").Collection("records")

	// User-defined filter string (from UI or config)
	userFilter := "A > 10 AND B < 5"
	parsedFilter := parseFilter(userFilter)

	// Add filter to only use recent 20 records
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(20)

	// Example: device-specific filter
	finalFilter := bson.D{
		{Key: "deviceId", Value: "abc-123"},
	}
	finalFilter = append(finalFilter, parsedFilter...)

	cursor, err := collection.Find(ctx, finalFilter, opts)
	if err != nil {
		panic(err)
	}

	var results []Record
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	fmt.Printf("Matching records: %d\n", len(results))
	for _, r := range results {
		fmt.Printf("A: %.2f, B: %.2f, Time: %s\n", r.A, r.B, r.CreatedAt.Format(time.RFC3339))
	}
}
