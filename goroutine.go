package main

import (
	"fmt"
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 완료 시 카운터 감소
	fmt.Printf("Worker %d 시작\n", id)
	// 실제로는 여기서 무거운 작업 수행
	fmt.Printf("Worker %d 완료\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg) // go 키워드 하나로 고루틴 시작
	}

	wg.Wait() // 모든 고루틴 완료 대기
	fmt.Println("모든 작업 완료")
}
