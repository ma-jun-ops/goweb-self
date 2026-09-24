package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ArticleCollection *mongo.Collection
var FileBucket *mongo.GridFSBucket

func ConnectDB() {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("mongodb连接失败", err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Ping失败", err)
	}
	db := client.Database("article")
	ArticleCollection = db.Collection("article")
	FileBucket = db.GridFSBucket()
	ArticleCollection = client.Database("article").Collection("articles")
	log.Println("mongodb连接成功")
}
