package main

import (
	"fmt"
	"runtime"
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 완료 시 카운터 감소
	fmt.Printf("Worker %d 시작\n", id)
	// 실제로는 여기서 무거운 작업 수행
	fmt.Printf("Worker %d 완료\n", id)
}

func main() { // 고루틴 실행중, 메인 함수가 끝나면 모든 고루틴도 강제 종료됨.
	var wg sync.WaitGroup
	fmt.Printf("초기 고루틴 수: %d\n", runtime.NumGoroutine())

	//Go 1.25 부터는 wg.go(함수명(파라미터)) 로 하면 Done,add 안해도 자동으로 추가됨.

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg) // go 키워드 하나로 고루틴 시작
	}
	fmt.Printf("생성 후 고루틴 수: %d\n", runtime.NumGoroutine())

	wg.Wait() // 모든 고루틴 완료 대기
	fmt.Println("모든 작업 완료")
	fmt.Printf("완료 후 고루틴 수: %d\n", runtime.NumGoroutine())
}
