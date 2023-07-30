package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017,localhost:27018,localhost:27019/my-database?replicaSet=my-replica-set"

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// var result bson.M
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
	// 	panic(err)
	// }
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	collection := client.Database("users").Collection("users")
	collection.InsertOne(context.TODO(), bson.D{
		{Key: "name", Value: "John Doe"},
		{Key: "age", Value: 43},
		{Key: "profession", Value: "Lawyer"},
	})

	var results []bson.M

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	if cur.All(context.Background(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}

}
