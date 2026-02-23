package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
메서드				데이터 소스	태그
ShouldBindJSON()	JSON 본문	json
ShouldBindQuery()	쿼리 스트링	form
ShouldBindUri()		URI 파라미터	uri
ShouldBindHeader()	HTTP 헤더	header
ShouldBind()		Content-Type 기반 자동 선택	form, json
*/

func main() {
	r := gin.Default()

	type CreateUserRequest1 struct { // JSON
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	r.POST("/users", func(c *gin.Context) { // JSON
		var req CreateUserRequest1

		// JSON 바인딩
		if err := c.ShouldBindJSON(&req); err != nil { //ShouldBindJSON 는 에러를 반환함.
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		})
	})

	//==================================================================================================================

	type LoginRequest struct { // HTML form
		Username string `form:"username"`
		Password string `form:"password"`
	}

	r.POST("/login", func(c *gin.Context) { // Form
		var req LoginRequest

		if err := c.ShouldBind(&req); err != nil { //ShouldBind 은 Content-Type 헤더 보고 자동으로 바인딩 해줌 application/json → JSON 바인딩 /application/x-www-form-urlencoded → 폼 바인딩  /multipart/form-data → 멀티파트 폼 바인딩
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username": req.Username,
		})
	})

	//==================================================================================================================

	type SearchRequest struct { // 쿼리 스트링을 구조체로 받을 수 있음
		Keyword  string `form:"q"`
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
	}

	r.GET("/search", func(c *gin.Context) {
		var req SearchRequest

		if err := c.ShouldBindQuery(&req); err != nil { //ShouldBindQuery 는 form 으로 지정한 문자 뒤에 쿼리파라미터의 값을 가져옴, ex) /search?q=golang&page=1&page_size=10 -> {"keyword":"golang","page":1,"page_size":10}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"keyword":   req.Keyword,
			"page":      req.Page,
			"page_size": req.PageSize,
		})
	})

	//==================================================================================================================

	type GetUserRequest struct {
		ID int `uri:"id" binding:"required"`
	}

	r.GET("/users/:id", func(c *gin.Context) {
		var req GetUserRequest

		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": req.ID,
		})
	})

	//==================================================================================================================

	type AuthHeader struct { // 헤더
		Authorization string `header:"Authorization"`
		UserAgent     string `header:"User-Agent"`
	}

	r.GET("/headers", func(c *gin.Context) {
		var h AuthHeader

		if err := c.ShouldBindHeader(&h); err != nil { // ShouldBindHeader 헤더값 가져옴
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"authorization": h.Authorization,
			"user_agent":    h.UserAgent,
		})
	})

	//==================================================================================================================

	type UpdateUserRequest struct { // 복합 바인딩
		ID    int    `uri:"id" binding:"required"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	r.PUT("/users/:id", func(c *gin.Context) {
		var req UpdateUserRequest

		// URI 바인딩
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// JSON 바인딩
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    req.ID,
			"name":  req.Name,
			"email": req.Email,
		})
	})
	r.Run(":8080")
}
