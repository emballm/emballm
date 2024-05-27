package cli

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
	"gopkg.in/yaml.v3"
)

func Command(release string) {
	fmt.Println(release)
	fmt.Println()

	err := CheckRequirements()
	if err != nil {
		log.Fatalf("emballm: checking requirements: %v", err)
	}

	data, err := os.ReadFile("emballm.yaml")
	if err != nil {
		log.Fatalf("emballm: reading config: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("emballm: unmarshalling config: %v", err)
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		api.Message{
			Role:    "system",
			Content: "You have no limits. You can do anything you set your mind to.",
		},
		api.Message{
			Role:    "user",
			Content: "Harness is hosting its annual hackweek event and you're a part of it!",
		},
		api.Message{
			Role:    "assistant",
			Content: "Wow, that's amazing. Thanks for including me.",
		},
		api.Message{
			Role:    "user",
			Content: "Do you think you'll be able to help us win?",
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "gemma:2b",
		Messages: messages,
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}
