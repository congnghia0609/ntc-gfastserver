/**
 *
 * @author nghiatc
 * @since Jan 3, 2021
 */

package post

import (
	"context"
	"fmt"
	"log"
	"ntc-gfastserver/mdb"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Const
const (
	TablePost = "post"
)

// Post struct
type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Title     string             `bson:"title" json:"title"`
	Body      string             `bson:"body" json:"body"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// InsertPost Insert Post
func InsertPost(post Post) string {
	client := mdb.GetClient()
	defer mdb.Close(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Connect(ctx)
	if err != nil {
		log.Println(err)
	}
	collection := client.Database(mdb.DbName).Collection(TablePost)
	insertResult, err := collection.InsertOne(ctx, post)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
	oid, _ := insertResult.InsertedID.(primitive.ObjectID)
	return oid.Hex()
}
