package main

import (
	"fmt"
	"math"
)

// 인터페이스 정의
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 구조체
type Circle struct {
	Radius float64
}

// struct{} - 이건 메모리를 사용하지 않는 빈 구조체. 채널의 신호전달용으로 사용됨.

type Person struct {
	Name string
	Age  int
}

type Address struct {
	City    string
	ZipCode string
}

type Employee struct {
	Person  // 익명 필드 (임베딩)
	Address // 익명 필드
	Company string
}

type Rectangle struct {
	Width, Height float64
}

// 메서드 구현 (포인터 리시버 - 구조체의 주소를 불러옴 * 안붙이면 구조체 값을 복사해옴)
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 인터페이스를 파라미터로 받는 함수
func printShape(s Shape) {
	fmt.Printf("넓이: %.2f, 둘레: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
	e := Employee{
		Person:  Person{Name: "Alice", Age: 30},
		Address: Address{City: "Seoul", ZipCode: "12345"},
		Company: "TechCorp",
	}

	// 임베딩된 필드에 직접 접근
	fmt.Println(e.Name)    // Alice
	fmt.Println(e.City)    // Seoul
	fmt.Println(e.Company) // TechCorp

	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 4, Height: 6},
	}
	for _, s := range shapes {
		printShape(s)
	}
}
