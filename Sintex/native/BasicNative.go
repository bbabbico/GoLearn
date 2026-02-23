package native

/*
#include <stdio.h>
#include <stdlib.h>

// C 함수: 두 수 더하기
int add(int a, int b) {
return a + b;
}

// C 함수: 문자열 출력
void hello(const char* name) {
printf("Hello from C, %s!\n", name);
fflush(stdout);
}
*/
import "C"
import "unsafe"

// Add Go 래퍼: main.go에서 사용하려면 Go 레퍼로 감싸서 실행
func Add(a, b int) int {
	return int(C.add(C.int(a), C.int(b)))
}

// Hello Go 래퍼: Go string -> C string 변환/해제까지 여기서 처리
func Hello(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.hello(cname)
}
