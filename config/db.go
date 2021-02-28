package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// database
var DB *mongo.Database

// collection
var Coll *mongo.Collection

func init() {
	// get a mongo sessions
	// Set connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	err = client.Disconnect(context.TODO())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB.")

	// Get db and collection as ref
	DB = client.Database("employee")
	Coll = DB.Collection("empSalary")
}
