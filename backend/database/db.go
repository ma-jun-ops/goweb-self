package database

import (
	"context"
	"log"
	"os"
	"regexp"
	"web/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ArticleCollection *mongo.Collection
var FileBucket *mongo.GridFSBucket

// envOrDefault 读取环境变量，未设置时返回默认值（本地开发无需配置）
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ConnectDB() {
	uri := envOrDefault("MONGO_URI", "mongodb://localhost:27017")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
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

// SearchArticles 关键词搜索：直接查 MongoDB（已停用 Elasticsearch）
func SearchArticles(keyword string) ([]model.Article, error) {
	ctx := context.Background()
	// QuoteMeta 转义正则特殊字符，避免用户输入导致正则报错
	pattern := bson.M{"$regex": regexp.QuoteMeta(keyword), "$options": "i"}
	filter := bson.M{"$or": []bson.M{
		{"title": pattern},
		{"content": pattern},
		{"tags": pattern},
	}}
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(50)

	cursor, err := ArticleCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	articles := make([]model.Article, 0)
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}
