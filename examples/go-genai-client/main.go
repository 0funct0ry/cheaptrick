package main

import (
	"context"
	"fmt"
	"log"

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

	fmt.Println("Sending prompt to local Gemini Mock Server (http://localhost:8080)...")
	prompt := "Tell me a short joke about a programmer."

	fmt.Printf("Prompt: %q\n", prompt)
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(prompt), nil)
	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	fmt.Println("\033[32mMock Server Response:\033[0m")
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		fmt.Printf("%s\n", resp.Text())
	} else {
		fmt.Printf("%+v\n", resp)
	}
}
