package Pattern

import (
	"fmt"
	"time"
)

func worker1(done <-chan struct{}) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			fmt.Println("종료 신호 수신")
			return
		case t := <-ticker.C:
			fmt.Println("작업:", t.Format("15:04:05"))
		}
	}
}

func main() {
	done := make(chan struct{})

	go worker1(done)

	time.Sleep(1 * time.Second)

	close(done) // 종료 신호
	time.Sleep(100 * time.Millisecond)
}
