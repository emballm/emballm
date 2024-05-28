package ollama

import (
	"context"
	"embed"
	"fmt"

	"github.com/ollama/ollama/api"
	"gopkg.in/yaml.v3"

	"emballm/internal/services"
)

//go:embed prompt.yaml
var content embed.FS

func Scan(model string, filePaths []string) (result *string, err error) {
	var prompt services.Prompt
	data, err := content.ReadFile("prompt.yaml")
	if err != nil {
		return nil, fmt.Errorf("emballm: reading prompt: %v", err)
	}
	err = yaml.Unmarshal(data, &prompt)
	if err != nil {
		return nil, fmt.Errorf("emballm: unmarshalling prompt: %v", err)
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("emballm: creating ollama client: %v", err)
	}

	var messages []api.Message
	for _, message := range prompt.Messages {
		messages = append(messages, api.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
	}

	respond := func(resp api.ChatResponse) error {
		if result == nil {
			stream := resp.Message.Content
			result = &stream
			return nil
		}
		stream := *result + resp.Message.Content
		result = &stream
		return nil
	}

	err = client.Chat(ctx, req, respond)
	if err != nil {
		return nil, fmt.Errorf("emballm: chat with ollama: %v", err)
	}

	return
}
