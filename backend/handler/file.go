package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"web/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有收到文件"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//存到GridFS，返回文件ID
	fileID, err := database.FileBucket.UploadFromStream(
		c.Request.Context(),
		fileHeader.Filename,
		bytes.NewReader(data),
		options.GridFSUpload(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fileId":   fileID.Hex(),
		"filename": fileHeader.Filename,
		"text":     extractDocxWithImages(c.Request.Context(), data),
	})
}

func DownloadFile(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效密码"})
		return
	}
	c.Header("Content-Disposition", "attachment")
	c.Header("Content-Type", "application/octet-stream")
	if _, err := database.FileBucket.DownloadToStream(c.Request.Context(), objID, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func ImageFile(c *gin.Context) {
	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var fileDoc struct {
		Filename string `bson:"filename"`
	}
	database.FileBucket.GetFilesCollection().FindOne(c.Request.Context(), bson.M{"id": objID}).Decode(&fileDoc)
	contentType := mime.TypeByExtension(filepath.Ext(fileDoc.Filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline")
	if _, err := database.FileBucket.DownloadToStream(c.Request.Context(), objID, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// 解析docx，返回文字图片占位
func extractDocxWithImages(ctx context.Context, data []byte) string {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return ""
	}
	var documentXML []byte
	rels := map[string]string{}
	media := map[string][]byte{}
	for _, file := range reader.File {
		switch {
		case file.Name == "word/document.xml":

			documentXML = readZipFile(file)
		case file.Name == "word/_rels/document.xml.rels":
			rels = parseRels(readZipFile(file))
		case strings.HasPrefix(file.Name, "word/media/"):
			media[file.Name] = readZipFile(file)
		}
	}
	text, rids := extractTextAndImages(documentXML)
	for _, rid := range rids {
		target := rels[rid]
		if target == "" {
			continue
		}
		imgBytes := media["word/"+target]
		if imgBytes == nil {
			imgBytes = media[target]
		}
		if imgBytes == nil {
			continue
		}
		imgID, err := database.FileBucket.UploadFromStream(
			ctx,
			filepath.Base(target),
			bytes.NewReader(imgBytes),
			options.GridFSUpload(),
		)
		if err != nil {
			continue
		}
		text = strings.Replace(text, "[IMG:"+rid+"]", "[IMG:"+imgID.String()+"]", 1)
	}
	return text
}

func extractTextAndImages(documentXML []byte) (string, []string) {
	type textNode struct {
		Text string `xml:",chardata"`
	}
	decoder := xml.NewDecoder(bytes.NewReader(documentXML))
	var builder strings.Builder
	var rids []string
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "t":
			var node textNode
			if err := decoder.DecodeElement(&node, &start); err == nil {
				builder.WriteString(node.Text)
			}
		case "blip":
			for _, attr := range start.Attr {
				if attr.Name.Local == "embed" {
					builder.WriteString("[IMG:" + attr.Value + "]")
					rids = append(rids, attr.Value)
				}
			}
		}
	}
	return builder.String(), rids
}

func parseRels(data []byte) map[string]string {
	rels := map[string]string{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}
		var id, target string
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "Id":
				id = attr.Value
			case "Target":
				target = attr.Value
			}
		}
		if id != "" && target != "" {
			rels[id] = target
		}
	}
	return rels
}

func readZipFile(file *zip.File) []byte {
	rc, err := file.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil
	}
	return data
}
