package main

import "fmt"

func classify(n int) string {
	// if 초기화 구문 지원
	if rem := n % 2; rem == 0 {
		return "짝수"
	} else {
		return "홀수"
	}
}

func main() {
	// Go에는 while이 없다 — for로 통일
	for i := 0; i < 5; i++ {
		fmt.Printf("%d는 %s\n", i, classify(i))
	}

	// while처럼 사용
	n := 1
	for n < 128 {
		n *= 2
	}

	// range로 슬라이스 순회
	fruits := []string{"🍎", "🍊", "🍋"}
	for idx, fruit := range fruits { // 배열 하나만 입력하면 index 값과 value가 자동으로 할당됨. 한개 배열, 한개 변수 입력하면 배열의 value 만 할당.
		fmt.Printf("[%d] %s\n", idx, fruit)
	}

	// switch — fallthrough 없음이 기본
	day := "토요일"
	switch day {
	case "토요일", "일요일":
		fmt.Println("주말!")
		fallthrough //원래 Go는 한번 case 만족하면 switch 바로 종료됨 이거 있으면 다음 case 문도 순회함
	default:
		fmt.Println("평일")

	}
}
