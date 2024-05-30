package vertex

import (
	"context"
	"emballm/internal/scans"
	"emballm/internal/scans/results"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"

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
			fileResult, err := ScanFile(scanClient, fileScan.Path)
			fileScan.Status = scans.Status.Complete
			if err != nil {
				fileScan.Status = scans.Status.Nope
				return
			}
			fileScan.Status = scans.Status.Complete

			// Create an instance of the Vulnerability struct
			result := strings.ReplaceAll(*fileResult, "```", "")
			result = strings.ReplaceAll(result, "json", "")

			issue := &results.Issue{}
			err = json.Unmarshal([]byte(result), issue)
			if err != nil {
				fileScan.Status = scans.Status.Nope
				return
			}
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

	gemini := client.GenerativeModel(scanClient.Model)
	chat := gemini.StartChat()

	var messages []*genai.Content

	for _, message := range prompt.Messages[:len(prompt.Messages)-1] {
		messages = append(messages, &genai.Content{
			Role:  message.Role,
			Parts: []genai.Part{genai.Text(message.Content)},
		})
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("vertex scan: reading file: %v", err)
	}
	codeMessage := prompt.Messages[len(prompt.Messages)-1]
	chat.History = append(chat.History, messages...)

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
