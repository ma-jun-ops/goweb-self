package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"web/database"

	"github.com/gin-gonic/gin"
)

type githubFriend struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	Bio       string `json:"bio"`
	HTMLURL   string `json:"htmlUrl"`
}

type githubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
	HTMLURL   string `json:"html_url"`
}

var githubHTTPClient = &http.Client{Timeout: 10 * time.Second}

func githubGet(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "my-little-yard")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := githubHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("github api 返回 %s", res.Status)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// GetGithubFriends 获取我关注的人
func GetGithubFriends(c *gin.Context) {
	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		username = "ma-jun-ops"
	}
	hasToken := os.Getenv("GITHUB_TOKEN") != ""

	const cacheKey = "github:following"
	if cached, err := database.RedisClient.Get(c.Request.Context(), cacheKey).Result(); err == nil {
		var friends []githubFriend
		if json.Unmarshal([]byte(cached), &friends) == nil {
			c.JSON(http.StatusOK, gin.H{"data": friends})
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancel()

	var following []githubUser
	url := fmt.Sprintf("https://api.github.com/users/%s/following?per_page=100", username)
	if err := githubGet(ctx, url, &following); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "获取 GitHub 关注列表失败"})
		return
	}

	// 未配置 Token 时 GitHub 限额为 60 次/小时，留出余量只取前 50 人
	limit := len(following)
	if !hasToken && limit > 50 {
		limit = 50
	}

	friends := make([]githubFriend, limit)
	sem := make(chan struct{}, 6)
	var wg sync.WaitGroup
	for i := 0; i < limit; i++ {
		u := following[i]
		wg.Add(1)
		go func(index int, user githubUser) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			friend := githubFriend{
				Login:     user.Login,
				Name:      user.Login,
				AvatarURL: user.AvatarURL,
				HTMLURL:   user.HTMLURL,
			}
			// 列表接口不含简介，再取一次用户详情；失败时保留基础信息
			var detail githubUser
			if err := githubGet(ctx, "https://api.github.com/users/"+user.Login, &detail); err == nil {
				if detail.Name != "" {
					friend.Name = detail.Name
				}
				friend.Bio = detail.Bio
				if detail.HTMLURL != "" {
					friend.HTMLURL = detail.HTMLURL
				}
				if detail.AvatarURL != "" {
					friend.AvatarURL = detail.AvatarURL
				}
			}
			friends[index] = friend
		}(i, u)
	}
	wg.Wait()

	if data, err := json.Marshal(friends); err == nil {
		database.RedisClient.Set(ctx, cacheKey, data, time.Hour)
	}
	c.JSON(http.StatusOK, gin.H{"data": friends})
}
