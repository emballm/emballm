package vertex

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"gopkg.in/yaml.v3"

	"emballm/internal/services"
)

//go:embed prompt.yaml
var content embed.FS

func Scan(model string, filePaths []string) (result *string, err error) {
	var prompt services.Prompt
	data, err := content.ReadFile("prompt.yaml")
	if err != nil {
		log.Fatalf("emballm: reading prompt: %v", err)
	}
	err = yaml.Unmarshal(data, &prompt)
	if err != nil {
		log.Fatalf("emballm: unmarshalling prompt: %v", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, "projectID", "location")
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
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

	r, err := chat.SendMessage(
		ctx,
		genai.Text(prompt.Messages[len(prompt.Messages)-1].Content),
	)
	if err != nil {
		return nil, err
	}
	rb, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json.MarshalIndent: %w", err)
	}
	resp := string(rb)
	result = &resp

	return
}
