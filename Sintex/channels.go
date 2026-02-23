package main

import (
	"fmt"
	"time"
)

// 생산자: 채널에 데이터 전송 - 언버퍼드
func generate_unbuffered(nums ...int) <-chan int {
	out := make(chan int) // 버퍼없이 송신 받으면 수신 받을때까지 블로킹됨.
	go func() {
		for _, n := range nums {
			out <- n // 채널로 전송
		}
		close(out) // 전송 완료 시 닫기
	}()
	return out
}

// 생산자: 채널에 데이터 전송 - 버퍼드
func generate_buffered(nums ...int) <-chan int {
	out := make(chan int, 5) //버퍼 5만큼 만들어서 5개 송신 까지는 수신자 안기다림. 버퍼 가득차면 블로킹됨.
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
	c1 := generate_unbuffered(2, 3, 4, 5)
	c2 := generate_unbuffered(2, 3, 4, 5)

	c3 := square(c1)
	for n := range c3 {
		fmt.Println(n) // 4, 9, 16, 25
	}

	// 버퍼 확인
	fmt.Printf("저장된 값: %d\n", len(c2))
	fmt.Printf("버퍼 크기: %d\n", cap(c2))
	fmt.Printf("남은 공간: %d\n", cap(c2)-len(c2))

	// select — 복수 채널 동시 대기
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	ch1 <- "from ch1"
	ch2 <- "from ch2"

	select { // select 로 여러 패널패턴 구현가능 https://wikidocs.net/320800
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	case <-time.After(1 * time.Second): //time.After는 지정된 시간 후에 현재 시간을 보내는 채널을 반환.
		fmt.Println("타임아웃!")
	}
}
