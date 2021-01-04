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

	"go.mongodb.org/mongo-driver/bson"
)

// Const
const (
	TablePost = "post"
)

// Post struct
type Post struct {
	// ID        primitive.ObjectID `bson:"_id" json:"id"`
	ID        int64     `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Body      string    `bson:"body" json:"body"`
	CreatedAt time.Time `bson:"created_at" json:"-"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// InsertPost Insert Post
func InsertPost(post Post) error {
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	insertResult, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		fmt.Println("Error Inserted post with insertResult:", insertResult)
		log.Println(err)
		return err
	}
	return nil
}

// GetPost get post
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-overview
func GetPost(id int64) Post {
	var post Post
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{{"_id", id}}
	err := collection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		log.Println(err)
	}
	return post
}

// GetAllPost get all post
func GetAllPost() []Post {
	var posts []Post
	client := mdb.GetClient()
	defer mdb.Close(client)
	collection := client.Database(mdb.DbName).Collection(TablePost)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			log.Println(err)
		}
		// do something with result...
		posts = append(posts, post)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	return posts
}
