package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("OPENROUTER_API_KEY")

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	// JSON 모드는 스트리밍을 끄고 한 번에 받는 것이 파싱하기 좋습니다.
	req := openai.ChatCompletionRequest{
		Model: "stepfun/step-3.5-flash:free",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				// 🚨 중요: 프롬프트 내용에 'JSON'이라는 단어가 반드시 포함되어야 합니다.
				Content: "너는 텍스트 분석기야. 사용자의 입력에서 이름과 나이를 추출해서 JSON 형태로 반환해줘. 예시: {\"name\": \"홍길동\", \"age\": 30}",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "내 이름은 김제미나이고, 올해 25살이야. 잘 부탁해!",
			},
		},
		// ✨ 핵심: 응답 포맷을 JSON Object로 강제 설정
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("API 요청 에러: %v", err)
	}

	// AI의 응답은 완벽한 JSON 문자열로 나옵니다.
	jsonOutput := resp.Choices[0].Message.Content
	fmt.Println("AI가 생성한 JSON 데이터:")
	fmt.Println(jsonOutput)

	// 이 jsonOutput을 Go의 json.Unmarshal()을 사용해 구조체(Struct)로 변환하면
	// 백엔드 데이터베이스에 바로 저장할 수 있습니다!
}
