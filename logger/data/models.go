package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoService struct {
	Client *mongo.Client
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {

	return Models{
		LogEntry: LogEntry{},
	}
}

func (mongoService *MongoService) Insert(entry LogEntry) error {
	collection := mongoService.Client.Database("logs").Collection("logs")

	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (mongoService *MongoService) GetAll() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := mongoService.Client.Database("logs").Collection("logs")

	options := options2.Find()
	options.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Println("Error while getting all the logs", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error while decoding the item from Mongo", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (mongoService *MongoService) GetById(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := mongoService.Client.Database("logs").Collection("logs")

	documentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	options := options2.FindOne()
	err = collection.FindOne(ctx, bson.M{"_id": documentId}, options).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (mongoService *MongoService) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := mongoService.Client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (mongoService *MongoService) Update(entry *LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := mongoService.Client.Database("logs").Collection("logs")

	documentId, err := primitive.ObjectIDFromHex(entry.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": documentId},
		bson.D{
			{"$set", bson.D{
				{"name", entry.Name},
				{"data", entry.Data},
				{"updated_at", entry.UpdatedAt},
			}},
		})

	if err != nil {
		return nil, err
	}

	return result, nil
}
