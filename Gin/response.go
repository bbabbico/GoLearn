package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 */
func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ // c.JSON()의 첫 번째 인자는 HTTP 상태 코드, 두 번째 인자는 응답 데이터입니다. 데이터는 자동으로 JSON으로 직렬화
			"message": "Hello, World!",
		})
	})

	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type User struct {
		ID      int     `json:"id"`
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	r.GET("/users/:id", func(c *gin.Context) { // 복잡한 JSON 응답
		user := User{
			ID:   1,
			Name: "Alice",
			Address: Address{
				City:    "Seoul",
				Country: "Korea",
			},
		}

		c.JSON(http.StatusOK, user)
	})

	r.GET("/users", func(c *gin.Context) { //슬라이스 응답
		users := []User{
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

	//================================================================================================================== 표준 응답 처리

	// 사용자 목록
	r.GET("/users", func(c *gin.Context) {
		users := []User{
			{1, "Alice", Address{City: "Seoul", Country: "Korea"}},
			{2, "Bob", Address{City: "Seoul", Country: "Korea"}},
		}

		List(c, users, 1, 10, 2)
	})

	// 사용자 상세
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		if id == "0" {
			NotFound(c, "사용자를 찾을 수 없습니다")
			return
		}

		user := User{1, "Alice", Address{City: "Seoul", Country: "Korea"}}
		Success(c, user)
	})

	// 사용자 생성
	r.POST("/users", func(c *gin.Context) {
		var req struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			BadRequest(c, err.Error())
			return
		}

		user := User{1, req.Name, Address{City: "Seoul", Country: "Korea"}}
		Created(c, user)
	})

}

// ====================================================================================================================== 응답 표준화

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Success 성공 응답
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created 생성 성공 응답
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// List 목록 응답
func List(c *gin.Context, data any, page, perPage, total int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		},
	})
}

// Fail 에러 응답
func Fail(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}

// BadRequest 400 Bad Request
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound 404 Not Found
func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, "NOT_FOUND", message)
}

// InternalError 500 Internal Server Error
func InternalError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
