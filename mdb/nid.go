/**
 *
 * @author nghiatc
 * @since Jan 4, 2021
 */

package mdb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Const
const (
	TableNId = "nlidgen"
)

// NId is long id gen
type NId struct {
	ID  string `bson:"_id" json:"_id"`
	Seq int64  `bson:"seq" json:"seq"`
}

// Next generate a auto increment version ID for the given key
func Next(id string) (int64, error) {
	var nid NId
	client := GetClient()
	defer Close(client)
	collection := client.Database(DbName).Collection(TableNId)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$inc", bson.D{{"seq", 1}}}}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&nid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		} else {
			log.Println(err)
			return 0, err
		}
	}
	// log.Println("nid:", nid)

	return nid.Seq + 1, nil
}
