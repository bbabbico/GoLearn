package main

import (
	"fmt"
	"os"
)

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("파일 열기 실패: %w", err) // %w로 에러 래핑
	}
	defer f.Close() // 함수 종료 시 자동 실행 — 리소스 정리에 필수

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func main() {
	// defer는 LIFO 순서로 실행
	defer fmt.Println("3: 마지막")
	defer fmt.Println("2: 중간")
	defer fmt.Println("1: 처음")
	fmt.Println("실행 중...")

	content, err := readFile("test.txt")
	if err != nil {
		fmt.Println("에러:", err)
		return
	}
	fmt.Println(content)
}
