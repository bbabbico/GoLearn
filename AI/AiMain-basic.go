package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {

	// 1. .env 파일 로드
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env 파일을 불러오는 데 실패했습니다.")
	}

	// 2. 환경 변수에서 API 키 가져오기
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY가 설정되지 않았습니다.")
	}

	//fmt.Printf(apiKey) // API 키 출력

	// 3. OpenRouter 설정 (BaseURL 변경)
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	// 4. 스트리밍 요청 설정
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: "stepfun/step-3.5-flash:free", // 사용할 OpenRouter 모델
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Go 언어로 백엔드를 개발할 때의 장점을 3가지 알려줘. (조금 길게 설명해줘)",
			},
		},
		Stream: true, // ✨ 스트리밍 활성화
	}

	// 5. 스트리밍 클라이언트 생성
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Fatalf("스트리밍 요청 에러: %v\n", err)
	}
	defer stream.Close()

	fmt.Print("AI 응답: \n\n")

	// 6. 실시간 응답 출력 루프 - 터미널 출력
	for {
		response, err := stream.Recv() // gRPC 스트리밍 통신에서 서버 또는 클라이언트로부터 다음데이터를 수신(Receive) 하는 표준 메서드

		// 응답이 끝났을 때 (End Of File)
		if errors.Is(err, io.EOF) {
			fmt.Println("\n\n[출력 완료]")
			break
		}

		// 수신 중 에러가 발생했을 때
		if err != nil {
			log.Printf("\n스트림 수신 중 에러 발생: %v\n", err)
			break
		}

		// ✨ 스트리밍은 Message.Content가 아니라 Delta.Content에 데이터가 담깁니다.
		fmt.Print(response.Choices[0].Delta.Content)
	}
}
