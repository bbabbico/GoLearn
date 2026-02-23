package main

import "fmt"

func main() {
	// var 키워드로 명시적 선언
	var name string = "Gopher"
	var age int = 13
	var pi float64 = 3.14159
	var ok bool = true
	var vv = 1 // 타입 추론

	// := 단축 선언 (함수 내부에서만 사용 가능)
	city := "Seoul"
	score := 98.5

	// 다중 선언
	x, y, z := 1, 2, 3

	//암묵적 type conversion 그딴거 없음.
	var i int = 100
	var u uint = uint(i)
	var f float32 = float32(i)
	println(f, u)

	str := "ABC"
	bytes := []byte(str)
	str2 := string(bytes)
	println(bytes, str2)

	// 포인터
	var p *int
	var k int = 10

	var _ = &k  //k의 주소를 할당
	println(*p) //p가 가리키는 주소에 있는 실제 내용을 출력

	// 상수
	const MaxSize = 1024
	const Pi float64 = 3.14159265358979
	const (
		Sky   = "Blue"
		Rose  = "Red"
		Gyuri = "Awesome"
	)

	const (
		apple = iota + 1 // 순차적으로 정수를 0부터 대입해줌
		grape
		orange
	)

	fmt.Printf(`\n\n\n\n\n 여기는 모든것이 문자여 문자가 해석되지 않음 %d`, vv)
	fmt.Printf("%s(%d)는 %s에 삽니다. score=%.1f + %d\n", name, age, city, score, vv)
	_ = pi
	_ = ok
	_ = x
	_ = y
	_ = z
	_ = MaxSize
}
