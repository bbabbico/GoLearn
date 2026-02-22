package main

import "fmt"

// 생산자: 채널에 데이터 전송
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n // 채널로 전송
		}
		close(out) // 전송 완료 시 닫기
	}()
	return out
}

// 소비자: 채널에서 데이터 수신하여 제곱
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in { // 채널이 닫힐 때까지 수신
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// 파이프라인 패턴
	c1 := generate(2, 3, 4, 5)
	c2 := square(c1)
	for n := range c2 {
		fmt.Println(n) // 4, 9, 16, 25
	}

	// select — 복수 채널 동시 대기
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	ch1 <- "from ch1"
	ch2 <- "from ch2"

	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	}
}
