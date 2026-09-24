package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Article struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Content   string        `json:"content" bson:"content"`
	Tags      []string      `json:"tags" bson:"tags"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type Attachment struct {
	FileID   string `json:"fileId" bson:"fileId"`
	FileName string `json:"filename" bson:"filename"`
	Text     string `json:"text" bson:"text"`
}
