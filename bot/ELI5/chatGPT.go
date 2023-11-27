package ELI5

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"os"
)

var client *openai.Client

func ConnectToChatGPT() {
	client = openai.NewClient(os.Getenv("OPENAI_KEY"))
	if client == nil {
		log.Fatal().Msg("Couldn't connect to OpenAI.")
	} else {
		log.Info().Msg("ChatGPT Connection initialized.")
	}
}

const initialPrompt = " Act as InputBOT until I say the keyword \"Done\" followed by my prompt. Your role is simply to take input from the user in this ChatGPT console and store it in your local memory for use after I say we are \"done\". We will begin the process now your first input is: "

func createChatCompletionMessage(content string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
}

func ProcessWords(data string) {
	everything := fmt.Sprintf("%s%s\n\n\nDone\n\nCould you explain the text in ~20 words like I am 5, please?", initialPrompt, data)
	var messages []openai.ChatCompletionMessage

	for i := 0; i < (len(everything)/4000)+1; i++ {
		messages = append(messages, createChatCompletionMessage(everything[i*4000:min(len(data)-1, 4000*(i+1))]))
	}

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})

	if err != nil {
		log.Error().Msg("Error with ChatGPT integration. Error: " + err.Error())
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
