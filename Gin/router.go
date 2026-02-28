package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

/*
// 성공
c.JSON(http.StatusOK, data)             // 200 - 조회 성공
c.JSON(http.StatusCreated, data)        // 201 - 생성 성공
c.JSON(http.StatusNoContent, nil)       // 204 - 삭제 성공 (본문 없음)

// 클라이언트 오류
c.JSON(http.StatusBadRequest, err)      // 400 - 잘못된 요청
c.JSON(http.StatusUnauthorized, err)    // 401 - 인증 필요
c.JSON(http.StatusForbidden, err)       // 403 - 권한 없음
c.JSON(http.StatusNotFound, err)        // 404 - 리소스 없음

// 서버 오류
c.JSON(http.StatusInternalServerError, err) // 500 - 서버 오류
*/

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"` // JSON 응답에서 제외
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"` // JSON 빈 값이면 제외
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type User1 struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

func main() {
	// Gin 라우터 생성
	r := gin.Default() //Default()는 기본 미들웨어인 Logger와 Recovery가 포함된 라우터를 반환.

	// 업로드 파일 크기 제한 (8MB)
	r.MaxMultipartMemory = 8 << 20

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

		c.JSON(http.StatusOK, gin.H{ //c.JSON 는 첫번째 인자는 HTTP 상태 코드, 두번쨰는 응답 본문. - gin.H는 map[string]any의 별칭 JSON 객체 편하게 만들 수 있게 해줌 map[string]any{"name": "Alice"} 이거랑 같음
			"message":    "Hello, Gin!", // 그냥 구조체로 받을 수도 있음 ex) c.JSON(http.StatusOK, gin.H{CreateUserRequest}
			"keyword":    keyword,
			"page":       page,
			"limit":      limit,
			"categories": categories,
			"config":     config,
		})
	})

	// 중첩 구조체
	r.GET("/users/:id", func(c *gin.Context) {
		user := User1{
			ID:   1,
			Name: "Alice",
			Address: Address{ //중첩 구조체
				City:    "Seoul",
				Country: "Korea",
			},
		}

		c.JSON(http.StatusOK, user)
	})
	/* 응답 결과
	{
	  "id": 1,
	  "name": "Alice",
	  "address": {
	    "city": "Seoul",
	    "country": "Korea"
	  }
	}
	*/

	type User2 struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// 목록 반환 슬라이스
	r.GET("/users", func(c *gin.Context) {
		users := []User2{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
			{ID: 3, Name: "Charlie"},
		}

		c.JSON(http.StatusOK, users)
	})
	/*
		[
		  {"id": 1, "name": "Alice"},
		  {"id": 2, "name": "Bob"},
		  {"id": 3, "name": "Charlie"}
		]

	*/
	r.POST("/upload", func(c *gin.Context) {
		// 폼에서 파일 가져오기
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 파일 저장
		filename := filepath.Base(file.Filename)
		dst := "./uploads/" + filename

		if err := c.SaveUploadedFile(file, dst); err != nil { // SaveUploadedFile : 파일 저장함수
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"filename": filename,
			"size":     file.Size,
		})
	})

	//배열 반환시 JSON 하이체킹 방지 응답은 JSON 배열 그대로 반환하면 털림. 객체로 감사서 "data" : [] 이런식으로 해야함
	r.GET("/data", func(c *gin.Context) {
		c.SecureJSON(http.StatusOK, []string{"a", "b", "c"})
	})

	// JSON HTML 이스케이프 후 반환
	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{"html": "<b>Hello</b>"})
	})

	// JSON 원본 반환
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{"html": "<b>Hello</b>"})
	})
	/*
		c.JSON()
		{"html":"\u003cb\u003eHello\u003c/b\u003e"} HTML 이스케이프 함

		c.PureJSON()
		{"html":"<b>Hello</b>"} HTML 이스케이프 안하고 원본 그대로 줌
	*/

	// 템플릿 로드
	r.LoadHTMLGlob("html/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{ // c.HTML 은 HTML 반환함. html 안의 {{ .title }} 요소로 변수 설정 가능
			"title": "홈페이지",
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
