package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chzyer/readline"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// Initialize the client configured to hit the mock server
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{
			BaseURL: "http://localhost:8080",
		},
		APIKey: "mock-key", // Any API key works since cheaptrick skips validation
	})
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[36mPrompt>\033[0m ",
		HistoryFile:     "/tmp/go-genai-client-history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatalf("Failed to initialize readline: %v", err)
	}
	defer rl.Close()

	fmt.Println("Connected to local Gemini Mock Server (http://localhost:8080)")
	fmt.Println("Enter your prompts below. Press Ctrl+D (EOF) to exit.")
	fmt.Println("---------------------------------------------------------")

	for {
		line, err := rl.Readline()
		if err != nil { // EOF (Ctrl+D) or Ctrl+C
			break
		}
		if line == "" {
			continue
		}

		// Generate content
		resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(line), nil)
		if err != nil {
			fmt.Printf("\033[31mError generating content: %v\033[0m\n", err)
			continue
		}

		// Print response text
		fmt.Println("\033[32mMock Server Response:\033[0m")
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			fmt.Printf("%s\n\n", resp.Text())
		} else {
			fmt.Printf("%+v\n\n", resp)
		}
	}
}
