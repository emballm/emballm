package ollama

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ollama/ollama/api"
	"gopkg.in/yaml.v3"

	"emballm/internal/scans"
	"emballm/internal/scans/results"
	"emballm/internal/services"
)

//go:embed prompt.yaml
var content embed.FS

func Scan(scanClient ScanClient, fileScans []*scans.FileScan) (issues []results.Issue, err error) {
	var scan []results.Issue
	var waitGroup sync.WaitGroup

	for _, fileScan := range fileScans {
		waitGroup.Add(1)
		go func(fileScan *scans.FileScan) {
			defer waitGroup.Done()
			fileResult, _ := ScanFile(scanClient, fileScan.Path)

			fileScan.Status = scans.Status.Complete

			// Create an instance of the Vulnerability struct
			result := strings.ReplaceAll(*fileResult, "```", "")
			result = strings.ReplaceAll(result, "json", "")

			issue := &results.Issue{}
			_ = json.Unmarshal([]byte(result), issue)
			issue.FileName = fileScan.Path

			scan = append(scan, *issue)
			fmt.Println(fmt.Sprintf("\t[%s] %s", scans.ScanStatus(fileScans), fileScan.Path))
		}(fileScan)
	}

	waitGroup.Wait()
	issues = scan
	return
}

func ScanFile(scanClient ScanClient, filePath string) (result *string, err error) {
	var prompt services.Prompt
	data, err := content.ReadFile("prompt.yaml")
	if err != nil {
		return nil, fmt.Errorf("ollama scan: reading prompt: %v", err)
	}
	err = yaml.Unmarshal(data, &prompt)
	if err != nil {
		return nil, fmt.Errorf("ollama scan: unmarshalling prompt: %v", err)
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("ollama scan: creating ollama client: %v", err)
	}

	var messages []api.Message
	for _, message := range prompt.Messages[:len(prompt.Messages)-1] {
		messages = append(messages, api.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ollama scan: reading file: %v", err)
	}
	codeMessage := prompt.Messages[len(prompt.Messages)-1]
	messages = append(messages, api.Message{
		Role:    codeMessage.Role,
		Content: fmt.Sprintf(codeMessage.Content, string(fileContent)),
	})

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    scanClient.Model,
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
		return nil, fmt.Errorf("ollama scan: chatting: %v", err)
	}

	return
}
