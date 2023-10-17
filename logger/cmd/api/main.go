package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger/data"
	"net/http"
	"time"
)

const (
	webPort  = "8085"
	rpcPost  = "5005"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models       data.Models
	MongoService data.MongoService
}

func main() {
	mongo, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongo

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close mongo connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models:       data.New(client),
		MongoService: data.MongoService{Client: mongo},
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("Logger is going to be served")
	err = server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	options := options2.Client().ApplyURI(mongoURL)
	options.SetAuth(options2.Credential{
		Username: "admin",
		Password: "password",
	})

	connection, err := mongo.Connect(context.TODO(), options)

	log.Println("Connected to Mongo")
	if err != nil {
		log.Println("Error connection to Mongo:", err)
	}

	return connection, nil
}
