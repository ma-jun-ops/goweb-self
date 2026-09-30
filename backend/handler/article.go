package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"web/database"
	"web/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateArticle(c *gin.Context) {
	var req struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Tags    []string `json:"tags"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	article := model.Article{
		ID:        bson.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := database.ArticleCollection.InsertOne(context.Background(), article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.PublishArticleEvent(database.ArticleEvent{
		Action:  "index",
		Article: &article,
		ID:      article.ID.Hex(),
	})
	c.JSON(http.StatusOK, gin.H{"message": "发布成功", "article": article})
}

func ListArticle(c *gin.Context) {
	ctx := context.Background()
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := database.ArticleCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	articles := []model.Article{}
	if err := cursor.All(ctx, &articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": articles})
}

func GetArticle(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}

	ctx := context.Background()
	cacheKey := "article:" + c.Param("id")

	// 1. 先查 Redis 缓存
	if cached, err := database.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		var article model.Article
		if json.Unmarshal([]byte(cached), &article) == nil {
			c.JSON(http.StatusOK, gin.H{"article": article})
			return
		}
	}

	// 2. 缓存未命中，查 MongoDB
	var article model.Article
	err = database.ArticleCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 写回缓存，TTL 5 分钟
	if data, err := json.Marshal(article); err == nil {
		database.RedisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func UpdateArticle(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}
	var req struct {
		Title   *string  `json:"title" `
		Content *string  `json:"content" `
		Tags    []string `json:"tags"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	setFields := bson.M{"updatedAt": time.Now()}
	if req.Title != nil {
		setFields["title"] = *req.Title
	}
	if req.Content != nil {
		setFields["content"] = *req.Content
	}
	if req.Tags != nil {
		setFields["tags"] = req.Tags
	}
	result, err := database.ArticleCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": setFields},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在 "})
		return
	}
	database.RedisClient.Del(context.Background(), "article:"+c.Param("id"))
	var updated model.Article
	if err := database.ArticleCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&updated); err != nil {
		database.PublishArticleEvent(database.ArticleEvent{
			Action:  "index",
			Article: &updated,
			ID:      updated.ID.Hex(),
		})
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章更新成功"})
}

func DeleteArticle(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}
	result, err := database.ArticleCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	database.RedisClient.Del(context.Background(), "article:"+c.Param("id"))
	database.PublishArticleEvent(database.ArticleEvent{
		Action: "delete",
		ID:     c.Param("id"),
	})
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func SearchArticle(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少搜索关键词"})
		return
	}
	articles, err := database.SearchArticles(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": articles})
}
