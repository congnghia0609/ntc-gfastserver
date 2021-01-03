/**
 *
 * @author nghiatc
 * @since Jan 3, 2021
 */

package mdb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Const
const (
	DbName = "fastdb"
)

// GetClient return mongo client
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Connect
func GetClient() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
		return nil
	}
	return client
}

// Close disconnect mongo client
func Close(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Println(err)
	}
}
