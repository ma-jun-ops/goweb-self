package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"web/model"

	"github.com/elastic/go-elasticsearch/v9"
)

const articleIndex = "articles"

var ESClient *elasticsearch.Client

func ConnectElasticsearch() {
	client, err := elasticsearch.New(
		elasticsearch.WithAddresses("http://localhost:9200"),
	)
	if err != nil {
		log.Fatal("elasticsearch连接失败", err)
	}
	res, err := client.Info()
	if err != nil {
		log.Fatal("elasticsearch连接失败", err)
	}
	res.Body.Close()
	ESClient = client
	ensureArticleIndex()
	log.Panicln("elasticsearch连接成功")
}

func ensureArticleIndex() {
	res, err := ESClient.Indices.Exists([]string{articleIndex})
	if err != nil {
		log.Fatal("检索索引失败", err)
	}
	exists := res.StatusCode == 200
	res.Body.Close()
	if exists {
		return
	}
	mapping := `{
		"mappings": {
			"properties": {
				"title":   {"type": "text"},
				"content": {"type": "text"},
				"tags":    {"type": "keyword"}
			}
		}
	}`
	res, err = ESClient.Indices.Create(articleIndex, ESClient.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Fatal("创建索引失败", err)
	}
	res.Body.Close()
}

// 把文章写入es
func IndexArticleDoc(article model.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	res, err := ESClient.Index(
		articleIndex,
		bytes.NewReader(data),
		ESClient.Index.WithDocumentID(article.ID.Hex()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("索引文档失败", res.String())
	}
	return nil
}

func DeleteArticleDoc(id string) error {
	res, err := ESClient.Delete(articleIndex, id)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("删除文档失败", res.String())
	}
	return nil
}

// 全文搜索
func SearchArticles(keyword string) ([]model.Article, error) {
	query := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  keyword,
				"fields": []string{"title", "content"},
			},
		},
	}
	queryBytes, err := json.Marshal(query)

	res, err := ESClient.Search(
		ESClient.Search.WithIndex(articleIndex),
		ESClient.Search.WithBody(bytes.NewReader(queryBytes)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("搜索失败: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source model.Article `json:"_source"`
			} `json:"hits"`
		}
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	articles := make([]model.Article, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		articles = append(articles, h.Source)
	}
	return articles, nil
}
