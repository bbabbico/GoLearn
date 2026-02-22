package main

import (
	"errors"
	"fmt"
)

// 다중 반환값 — Go의 핵심 패턴 func 함수명(파라미터) (반환타입) (반환 타입말고 반환 변수 명으로 설정가능)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("0으로 나눌 수 없습니다")
	}
	return a / b, nil
}

// 가변 인자 (variadic)
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 클로저
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 일급 함수 (함수를 파라미터로 받음)
func doWork(f func() int) {
	n := f()
	fmt.Println(n)
}

// type문을 이용해 두 정수를 합하는 함수형을 calculatorNum으로 정의
type calculatorNum func(int, int) int

// type문을 이용해 두 문자열을 복제하는 함수형을 calculatorStr로 정의
type calculatorStr func(string, string) string

func calNum(f calculatorNum, a int, b int) int {
	result := f(a, b)
	return result
}

func calStr(f calculatorStr, a string, b string) string {
	sentence := f(a, b)
	return sentence
}

func main() {
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("에러:", err)
	} else {
		fmt.Printf("결과: %.4f\n", result)
	}

	fmt.Println(sum(1, 2, 3, 4, 5)) // 15

	c := counter()             //익명함수
	fmt.Println(c(), c(), c()) // 1 2 3

	func(a, b int) int { // 이거를 파라미터로 넘기는것도 가능
		return a + b
	}(3, 5) // 정의 후 바로 호출

	multi := func(i int, j int) int {
		return i * j
	}
	duple := func(i string, j string) string {
		return i + j + i + j
	}

	r1 := calNum(multi, 10, 20)
	fmt.Println(r1)

	r2 := calStr(duple, "Hello", " Golang ")
	fmt.Println(r2)
}
