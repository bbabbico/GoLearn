package main

import "fmt"

func main() {
	// 배열 — 고정 길이
	var arr [3]int = [3]int{10, 20, 30}
	fmt.Println(arr[0]) // 10

	// 슬라이스 — 동적 배열, Go에서 가장 많이 쓰임
	nums := []int{1, 2, 3, 4, 5}
	nums = append(nums, 6) // 원소 추가
	sub := nums[1:4]       // 슬라이싱: [2 3 4]
	fmt.Println(sub)

	// make로 슬라이스 생성 (len=0, 최대 용량(cap=10)) 오직 slice , map , Channel 타입만 가능 반환값은 주소값이 아니라 변수 값 자체를 반환 -> []string 인 s를 반환하는게 아니라 {0,0,0,0,0} 이걸 반환
	s := make([]string, 0, 10)
	s = append(s, "Go", "Rust", "C")

	//new는 어떤타입이든 메모리 할당하고 주소값 반환. 값은 자동으로 0으로 초기화 하고
	var p *int = new(int)
	fmt.Println(*p) // 출력: 0, `new`에 의해 int에 할당된 메모리는 제로 값으로 초기화됨

	// 맵 — key-value 자료구조
	scores := map[string]int{
		"Alice": 95,
		"Bob":   82,
	}
	scores["Charlie"] = 91

	// 값 존재 여부 확인 (ok idiom)
	val, ok := scores["Dave"] //키 존재하면 ok 에 true 들어감
	if !ok {
		fmt.Println("Dave 없음, 기본값:", val)
	}
}
