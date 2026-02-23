package Pattern

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
}

func workerPool(numWorkers int, jobs <-chan Job, results chan<- Result) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				// 작업 처리
				output := fmt.Sprintf("워커%d가 %s 처리", workerID, job.Data)
				time.Sleep(100 * time.Millisecond)
				results <- Result{JobID: job.ID, Output: output}
			}
		}(i + 1)
	}

	wg.Wait()
	close(results)
}

func main() {
	jobs := make(chan Job, 10)
	results := make(chan Result, 10)

	// 워커 풀 시작
	go workerPool(3, jobs, results)

	// 작업 제출
	go func() {
		for i := 1; i <= 9; i++ {
			jobs <- Job{ID: i, Data: fmt.Sprintf("데이터%d", i)}
		}
		close(jobs)
	}()

	// 결과 수집
	for r := range results {
		fmt.Printf("작업 %d: %s\n", r.JobID, r.Output)
	}
}
