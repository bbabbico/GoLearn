package Pattern

import (
	"fmt"
	"sync"
)

func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func producer(id int, count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 1; i <= count; i++ {
			out <- id*100 + i
		}
		close(out)
	}()
	return out
}

func main() {
	ch1 := producer(1, 3) // 101, 102, 103
	ch2 := producer(2, 3) // 201, 202, 203
	ch3 := producer(3, 3) // 301, 302, 303

	merged := fanIn(ch1, ch2, ch3)

	for v := range merged {
		fmt.Println(v)
	}
}
