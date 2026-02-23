package Pattern

import (
	"errors"
	"fmt"
)

type Result1 struct {
	Value int
	Err   error
}

func divide(a, b int) <-chan Result1 {
	resultCh := make(chan Result1, 1)
	go func() {
		if b == 0 {
			resultCh <- Result1{Err: errors.New("0으로 나눌 수 없음")}
		} else {
			resultCh <- Result1{Value: a / b}
		}
		close(resultCh)
	}()
	return resultCh
}

func main() {
	r1 := <-divide(10, 2)
	if r1.Err != nil {
		fmt.Println("에러:", r1.Err)
	} else {
		fmt.Println("결과:", r1.Value)
	}

	r2 := <-divide(10, 0)
	if r2.Err != nil {
		fmt.Println("에러:", r2.Err)
	} else {
		fmt.Println("결과:", r2.Value)
	}
}
