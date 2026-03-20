package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// 1. 환경 변수 및 클라이언트 설정
	godotenv.Load()
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY가 설정되지 않았습니다.")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)
	ctx := context.Background()

	// 대화 기록을 최대 6개(질문 3번, 답변 3번)만 유지한다고 가정
	const maxHistory = 6

	// 2. 대화 기록을 저장할 슬라이스(배열) 생성
	// 첫 메시지로 AI의 역할을 부여할 수도 있음
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "너는 친절하고 명확하게 답변해주는 Go 언어 전문가야.",
		},
	}

	// 터미널 입력을 받기 위한 Reader 생성
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("🤖 챗봇이 시작되었습니다. (종료하려면 'q' 입력)")
	fmt.Println(strings.Repeat("-", 40))

	// 3. 무한 루프로 대화 진행
	for {
		// 사용자 입력 받기
		fmt.Print("\n👤 나: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if userInput == "q" {
			fmt.Println("채팅을 종료합니다.")
			break
		}
		if userInput == "" {
			continue
		}

		// ✨ 핵심 1: 사용자의 질문을 대화 기록에 추가
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		// 시스템 프롬프트(1개) + 최대 대화 기록(6개)을 초과했는지 확인
		if len(messages) > maxHistory+1 {
			// messages[0] : 시스템 프롬프트는 무조건 유지
			// messages[len(messages)-maxHistory:] : 가장 최근 10개의 대화만 가져오기

			// 두 개를 이어 붙여서 새로운 messages 배열로 덮어씌움
			messages = append([]openai.ChatCompletionMessage{messages[0]}, messages[len(messages)-maxHistory:]...)
		}

		// 요청 생성 (누적된 messages 배열 전체를 보냄)
		req := openai.ChatCompletionRequest{
			Model:    "stepfun/step-3.5-flash:free",
			Messages: messages,
			Stream:   true,
		}

		stream, err := client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			fmt.Printf("\n[오류] API 요청 실패: %v\n", err)
			continue
		}

		fmt.Print("🤖 AI: ")

		// AI의 답변을 하나로 합치기 위한 변수
		var assistantReply string

		// 스트리밍 수신 처리
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("\n[오류] 스트림 수신 실패: %v\n", err)
				break
			}

			// 화면에 출력하면서 동시에 답변을 변수에 저장
			chunk := response.Choices[0].Delta.Content
			fmt.Print(chunk)
			assistantReply += chunk
		}
		stream.Close()
		fmt.Println()

		// ✨ 핵심 2: 완성된 AI의 답변을 대화 기록에 추가 (다음 질문의 문맥이 됨)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantReply,
		})
	}
}
