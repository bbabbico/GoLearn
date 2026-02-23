package Pattern

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("워커 %d: 작업 %d 처리\n", id, job)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// 팬아웃: 3개 워커로 분배
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// 작업 전송
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
}
