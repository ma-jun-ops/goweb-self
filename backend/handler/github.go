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

// 返回给前端的github好友结构体
type githubFriend struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarurl"`
	Bio       string `json:"bio"`
	HTMLURL   string `json:"htmlurl"`
}

// github api返回的用户结构体
type githubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
	HTMLURL   string `json:"html_url"`
}

// 全局复用一个http.Client
var githubHttpClient = &http.Client{Timeout: 10 * time.Second}

// githubGet是对请求的统一封装
func githubGet(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	//请求必须带User-Agent
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "web-backend")

	//配置token
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := githubHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//检查响应状态码
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("github api返回错误状态码: %s", res.Status)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// githubGetFriends获取github好友
func GetGithubFriends(c *gin.Context) {
	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		username = "majun-ops"
	}

	hasToken := os.Getenv("GITHUB_TOKEN") != ""

	//先读redis
	const cacheKey = "github:following"
	if cached, err := database.RedisClient.Get(c.Request.Context(), cacheKey).Result(); err == nil {
		var friends []githubFriend
		if json.Unmarshal([]byte(cached), &friends) != nil {
			c.JSON(http.StatusOK, gin.H{"data": friends})
			return
		}

	}

	//创建一个20秒的总超时
	ctx, cancle := context.WithTimeout(c.Request.Context(), 20*time.Second)
	defer cancle()

	//拉取关注列表
	var following []githubUser
	url := fmt.Sprintf("https://api.github.com/users/%s/following?per_page=100", username)
	if err := githubGet(ctx, url, &following); err != nil {
		// 列表都拿不到，说明 GitHub 不可用或额度耗尽，直接返回 502
		c.JSON(http.StatusBadGateway, gin.H{"error": "获取 GitHub 关注列表失败"})
		return
	}

	//补token限制人数
	limit := len(following)
	if hasToken && limit > 50 {
		limit = 50
	}

	//并发补全个人信息
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
			var detail githubUser
			if err := githubGet(ctx, "https://api.github.com/users/"+user.Login, &detail); err != nil {
				if detail.Name != "" {
					friend.Name = detail.Name
				}
				friend.Bio = detail.Bio
			}
			if detail.HTMLURL != "" {
				friend.HTMLURL = detail.HTMLURL
			}
			if detail.AvatarURL != "" {
				friend.AvatarURL = detail.AvatarURL
			}
			friends[index] = friend
		}(i, u)

	}

	wg.Wait()

	// ---------- 第六步：写入缓存 ---------- // 缓存 1 小时。序列化失败时忽略错误，不影响本次返回。
	if data, err := json.Marshal(friends); err == nil {
		database.RedisClient.Set(ctx, cacheKey, data, time.Hour)
	}
	// ---------- 第七步：返回给前端 ---------- // 统一包一层 data，前端取 res.data.data
	c.JSON(http.StatusOK, gin.H{"data": friends})
}
