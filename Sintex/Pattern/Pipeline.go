package Pattern

import "fmt"

// 1단계: 숫자 생성
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// 2단계: 제곱
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 3단계: 문자열 변환
func stringify(in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		for n := range in {
			out <- fmt.Sprintf("결과: %d", n)
		}
		close(out)
	}()
	return out
}

func main() {
	// 파이프라인 연결
	nums := generator(1, 2, 3, 4, 5)
	squares := square(nums)
	results := stringify(squares)

	for r := range results {
		fmt.Println(r)
	}
}
