// filename: main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

/*
//문자열//
태그			설명			예시
required	필수			binding:"required"
min=N		최소 길이	binding:"min=2"
max=N		최대 길이	binding:"max=100"
len=N		정확한 길이	binding:"len=6"
email		이메일 형식	binding:"email"
url	URL 	형식			binding:"url"
alpha		알파벳만		binding:"alpha"
alphanum	알파벳+숫자만	binding:"alphanum"

//숫자//
태그		설명		예시
gt=N	초과		binding:"gt=0"
gte=N	이상		binding:"gte=1"
lt=N	미만		binding:"lt=100"
lte=N	이하		binding:"lte=99"
eq=N	같음		binding:"eq=10"
ne=N	같지 않음 binding:"ne=0"

//비교//
태그				설명				예시
oneof=a b c		지정된 값 중 하나	binding:"oneof=active inactive"
eqfield=Field	다른 필드와 같음	binding:"eqfield=Password" 비밀번호 확인란

omitempty - 값이 없으면 검증 건너뜀 ex) Age int `json:"age" binding:"omitempty,gte=0,lte=150"` 이렇게 되면 Age 값 없으면 그냥 건너뜀
*/

type CreateUserRequest2 struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gte=0,lte=150"`
	Password string `json:"password" binding:"required,min=8"`
}

func main() {
	r := gin.Default()

	r.POST("/users", func(c *gin.Context) {
		var req CreateUserRequest2

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		})
	})

	// 회원가입 예제
	type SignupRequest struct {
		Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=8"`
		PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
		Age             int    `json:"age" binding:"omitempty,gte=0,lte=150"`
		Role            string `json:"role" binding:"omitempty,oneof=user admin"`
	}

	r.POST("/signup", func(c *gin.Context) {
		var req SignupRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok { // 검증 에러 처리
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": formatValidationErrors(validationErrors),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 검증 에러 말고 다른에러 발생시
			return
		}

		c.JSON(http.StatusCreated, gin.H{"username": req.Username})
	})

	r.Run(":8080")
}

// 에러 메시지 커스텀
func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		tag := e.Tag()

		switch tag {
		case "required":
			errors[field] = field + "은(는) 필수입니다"
		case "email":
			errors[field] = "유효한 이메일 주소를 입력하세요"
		case "min":
			errors[field] = field + "은(는) 최소 " + e.Param() + "자 이상이어야 합니다"
		case "max":
			errors[field] = field + "은(는) 최대 " + e.Param() + "자까지 가능합니다"
		case "eqfield":
			errors[field] = field + "이(가) 일치하지 않습니다"
		default:
			errors[field] = field + " 값이 유효하지 않습니다"
		}
	}

	return errors
}
