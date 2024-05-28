package vertex

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"

	"emballm/internal/services"
)

//go:embed prompt.yaml
var content embed.FS

func Scan(model string, filePath string) (result *string, err error) {
	var prompt services.Prompt
	data, err := content.ReadFile("prompt.yaml")
	if err != nil {
		log.Fatalf("vertex scan: reading prompt: %v", err)
	}
	err = yaml.Unmarshal(data, &prompt)
	if err != nil {
		log.Fatalf("vertex scan: unmarshalling prompt: %v", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("VERTEX_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer func(client *genai.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("")
		}
	}(client)

	gemini := client.GenerativeModel(model)
	chat := gemini.StartChat()

	var messages []*genai.Content
	for _, message := range prompt.Messages[:len(prompt.Messages)-1] {
		message := genai.Content{
			Role:  message.Role,
			Parts: []genai.Part{genai.Text(message.Content)},
		}
		messages = append(messages, &message)
	}
	chat.History = append(chat.History, messages...)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("vertex scan: reading file: %v", err)
	}
	codeMessage := prompt.Messages[len(prompt.Messages)-1]
	r, err := chat.SendMessage(
		ctx,
		genai.Text(fmt.Sprintf(codeMessage.Content, string(fileContent))))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	rb, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	resp := string(rb)
	result = &resp

	return
}
