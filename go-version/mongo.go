package main

import (
	"context"
	"fmt"
	"iter"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	students *mongo.Collection
	client   *mongo.Client
}

func NewMongoDB() *MongoDB {
	client, err := mongo.Connect(options.Client().
		ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var result bson.M
	if err := client.
		Database("school").
		RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).
		Decode(&result); err != nil {
		panic(err)
	}

	return &MongoDB{
		students: client.Database("school").Collection("students"),
		client:   client,
	}
}

func (m *MongoDB) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		fmt.Printf("failed to disconnect from MongoDB, err: %s\n", err.Error())
	}
}

func (m *MongoDB) GetAllPagedData(itemsPerPage int64) iter.Seq[[]Student] {
	return func(yield func(students []Student) bool) {
		students := make([]Student, 0, itemsPerPage)
		var skip int64

		for {
			cursor, err := m.students.Find(context.Background(), bson.D{},
				options.Find().SetSkip(skip).SetLimit(itemsPerPage))
			if err != nil {
				fmt.Printf("failed to find students from skip %d, err: %s\n", skip, err.Error())
				return
			}

			if err := cursor.All(context.Background(), &students); err != nil {
				fmt.Printf("failed to decode students from skip %d, err: %s\n", skip, err.Error())
				return
			}

			if len(students) == 0 {
				return
			}

			if !yield(students) {
				return
			}

			skip += itemsPerPage
		}
	}
}

func (m *MongoDB) Count() (int, error) {
	count, err := m.students.CountDocuments(context.Background(), bson.D{}, options.Count().SetHint("_id_"))
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
