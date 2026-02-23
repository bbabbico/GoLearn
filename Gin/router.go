package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Gin 라우터 생성
	r := gin.Default() //Default()는 기본 미들웨어인 Logger와 Recovery가 포함된 라우터를 반환.

	// GET / 요청 처리
	r.GET("/g", func(c *gin.Context) { // GET 핸들러 함수, *gin.Context 가 요청 정보와 응답 매서드 가지고 있음.

		// 쿼리 파라미터 값
		keyword := c.Query("q")

		// 배열로 받기
		categories := c.QueryArray("category") // /g?category=a&category=b&category=c에 접속하면 categories는 ["a", "b", "c"]

		// 맵으로 받기
		config := c.QueryMap("config") //g?config[theme]=dark&config[lang]=ko에 접속하면 config는 {"theme": "dark", "lang": "ko"}

		// 기본값 지정
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		c.JSON(http.StatusOK, gin.H{ //c.JSON 는 첫번째 인자는 HTTP 상태 코드, 두번쨰는 응답 본문. - gin.H는 map[string]any의 별칭 JSON 객체 편하게 만들 수 있게 해줌
			"message":    "Hello, Gin!",
			"keyword":    keyword,
			"page":       page,
			"limit":      limit,
			"categories": categories,
			"config":     config,
		})
	})

	// GET 파라미터
	r.GET("/users/:userId/posts/:/*name", func(c *gin.Context) { // : 는 해당 위치의 값을 파라미터로 받음. * 는 그 위치 부터 뒤에 모든것을 파라미터로 받음. ex) /name/qwe/eeee에 접속하면 name은 /qwe/eeee
		userId := c.Param("userId")
		postId := c.Param("postId")

		c.JSON(http.StatusOK, gin.H{
			"user_id": userId,
			"post_id": postId,
		})
	})

	// POST - 생성
	r.POST("/users", func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil { // c.ShouldBindJSON() - 요청 본문의 JSON 을 구조체로 변환
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 여기서 실제로는 DB에 저장
		c.JSON(http.StatusCreated, gin.H{
			"id":    1,
			"name":  req.Name,
			"email": req.Email,
		})
	})

	// PUT - 전체 수정
	r.PUT("/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	})

	// DELETE - 삭제
	r.DELETE("/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})

	// 모든 HTTP 메서드 처리
	r.Any("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// 정의되지 않은 라우트 처리 (404)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
	})

	// /api/v1 그룹 생성
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		v1.Use(Middleware()) // v1의 모든 핸들러들은 Middleware() 함수를 거쳐서 실행됨
		{
			v1.GET("/users", func(c *gin.Context) { // /api/v1/users
				c.JSON(http.StatusOK, gin.H{"version": "v1", "users": []string{}})
			})
			v1.GET("/posts", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"version": "v1", "posts": []string{}})
			})
		}
		// v2 API (새로운 버전)
		v2 := r.Group("/api/v2")
		{
			v2.GET("/users", getUsersV2) // 응답 형식 변경 - 함수로 따로 빼서 사용가능
		}

	}

	// 서버 시작
	r.Run(":8080")
}

func getUsersV2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []map[string]any{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
		},
		"meta": gin.H{"total": 2},
	})
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
