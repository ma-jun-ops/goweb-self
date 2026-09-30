package database

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"web/model"

	"github.com/segmentio/kafka-go"
)

const articleTopic = "article-sync"

var KafkaWriter *kafka.Writer

// 发送到kafka的事件
type ArticleEvent struct {
	Action  string         `json:"action"`
	Article *model.Article `json:"article,omitempty"`
	ID      string         `json:"id"`
}

func ConnectKafka() {
	createTopic()
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    articleTopic,
		Balancer: &kafka.LeastBytes{},
	}
	go consumeArticles()
	log.Println("Kafka 连接成功")
}

func createTopic() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal("kafka连接失败", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("获取kafka controller失败", err)
	}
	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		log.Fatal("连接kafka controller失败", err)
	}
	defer controllerConn.Close()
	_ = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             articleTopic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

// 生产一条文章事件
func PublishArticleEvent(event ArticleEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return KafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	})
}

// 消费者：读事件写入Es
func consumeArticles() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   articleTopic,
		GroupID: "article-es-sync",
	})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("消费消息失败", err)
			continue
		}

		var event ArticleEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Fatal("解析事件失败", err)
			continue
		}

		switch event.Action {
		case "index":
			if event.Article != nil {
				if err := IndexArticleDoc(*event.Article); err != nil {
					log.Println("写入es失败", err)
				}
			}
		case "delete":
			if err := DeleteArticleDoc(event.ID); err != nil {
				log.Println("删除es文档失败", err)
			}
		}
	}
}
