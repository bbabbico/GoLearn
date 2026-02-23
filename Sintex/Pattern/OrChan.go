package Pattern

import (
	"fmt"
	"time"
)

func or(channels ...<-chan struct{}) <-chan struct{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan struct{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

func after(d time.Duration) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(done)
	}()
	return done
}

func main() {
	start := time.Now()

	<-or(
		after(2*time.Second),
		after(5*time.Second),
		after(1*time.Second), // 가장 빠름
		after(3*time.Second),
	)

	fmt.Printf("완료: %v\n", time.Since(start))
}
