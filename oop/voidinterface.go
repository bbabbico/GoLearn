package main

import "fmt"

func printAny(v any) { //any 는 빈 인터페이스의 별칭임.
	fmt.Printf("값: %v, 타입: %T\n", v, v)
}

func main() {
	var i interface{}

	i = 42
	fmt.Printf("값: %v, 타입: %T\n", i, i) // 값: 42, 타입: int

	i = "hello"
	fmt.Printf("값: %v, 타입: %T\n", i, i) // 값: hello, 타입: string

	i = true
	fmt.Printf("값: %v, 타입: %T\n", i, i) // 값: true, 타입: bool

	printAny(42)             // 값: 42, 타입: int
	printAny("hello")        // 값: hello, 타입: string
	printAny([]int{1, 2, 3}) // 값: [1 2 3], 타입: []int

	// 빈 인터페이스의 활용 1, 다양한 타입 저장 TODO: 빈 인터페이스는 타입 안정성 검사를 해야함
	items := []any{
		1,
		"hello",
		true,
		3.14,
		[]int{1, 2, 3},
	}

	for _, item := range items {
		fmt.Printf("%T: %v\n", item, item)
	}

	// 빈 인터페이스의 활용 2, 맵 값으로 사용
	data := map[string]any{
		"name":   "Alice",
		"age":    30,
		"active": true,
	}

	fmt.Println(data["name"]) // Alice

	// Go 제네릭
	nums := []int{1, 2, 3}
	fmt.Println(first(nums)) // 1

	strs := []string{"a", "b", "c"}
	fmt.Println(first(strs)) // a

	fmt.Println(Sum(10, 20))               // int 타입
	fmt.Println(Sum(10.5, 20.5))           // float64 타입
	fmt.Println(Sum(int32(10), int32(20))) // int32 타입

}

// 빈 인터페이스 대신 제네릭
func first[T any](slice []T) T { //any의 T 타입을 받는다는말, T는 뒤에 명시된 타입이나, 인터페이스와 인터페이스 안의 타입들을 받을 수 있음.
	return slice[0]
}

// Number 제네릭 타입 정의
type Number interface {
	~int | ~float64 | ~int32
}

// Sum 제네릭을 사용한 합계 계산 함수
func Sum[T Number](a, b T) T {
	return a + b
}
