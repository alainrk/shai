package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Config struct {
	LLMAPIURL string `mapstructure:"LLM_API_URL"`
	LLMAPIKey string `mapstructure:"LLM_API_KEY"`
	Model     string `mapstructure:"LLM_MODEL"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var (
	configPath string
	config     Config
	messages   []Message
	r          *glamour.TermRenderer
)

func init() {
	// Set up glamour renderer
	var err error
	r, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		fmt.Printf("Error setting up renderer: %v\n", err)
		os.Exit(1)
	}

	// Get default config path
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
		os.Exit(1)
	}
	defaultConfigPath := filepath.Join(usr.HomeDir, ".config", "shai", "config")

	// Parse command line flags
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to config file")
	flag.Parse()
}

func loadConfig() error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	if config.LLMAPIURL == "" {
		return fmt.Errorf("LLM_API_URL is required in config")
	}
	if config.LLMAPIKey == "" {
		return fmt.Errorf("LLM_API_KEY is required in config")
	}
	if config.Model == "" {
		return fmt.Errorf("LLM_MODEL is required in config")
	}

	return nil
}

func sendMessage(content string) (string, error) {
	// Add user message to history
	messages = append(messages, Message{
		Role:    "user",
		Content: content,
	})

	reqBody := ChatRequest{
		Model:    config.Model,
		Messages: messages,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", config.LLMAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LLMAPIKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s (%d)", string(body), resp.StatusCode)
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Check if there are choices
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	// Add assistant response to message history
	assistantMessage := chatResp.Choices[0].Message
	messages = append(messages, Message{
		Role:    "assistant",
		Content: assistantMessage.Content,
	})

	return assistantMessage.Content, nil
}

func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "exit" || input == "quit" {
		fmt.Println("Goodbye!")
		os.Exit(0)
	}

	if input == "" {
		return
	}

	fmt.Print("Thinking...")

	response, err := sendMessage(input)

	// Clear the thinking indicator
	fmt.Print("\r                \r")
	fmt.Print("\033[K")

	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	// Render markdown
	out, err := r.Render(response)
	if err != nil {
		color.Red("Error rendering response: %v", err)
		fmt.Println(response) // Fallback to plain text
		return
	}

	// Print rendered markdown
	fmt.Println("\n" + out)
}

func completer(d prompt.Document) []prompt.Suggest {
	// This could be enhanced to provide contextual suggestions
	return []prompt.Suggest{}
}

func main() {
	// Banner
	color.Cyan("\n=== Shai - LLM Shell ===\n")

	if err := loadConfig(); err != nil {
		if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
			// Create default config directory
			configDir := filepath.Dir(configPath)
			if err := os.MkdirAll(configDir, 0755); err != nil {
				color.Red("Error creating config directory: %v", err)
				os.Exit(1)
			}

			// Create default config file
			defaultConfig := `LLM_API_URL=https://api.openai.com/v1/chat/completions
LLM_API_KEY=your_api_key_here
LLM_MODEL=gpt-4
`
			if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
				color.Red("Error creating default config file: %v", err)
				os.Exit(1)
			}

			color.Yellow("Created default config at %s", configPath)
			color.Yellow("Please edit the file and set your API key before continuing.")
			os.Exit(1)
		}

		color.Red("Error loading config: %v", err)
		os.Exit(1)
	}

	color.Green("Connected to %s using model %s", config.LLMAPIURL, config.Model)
	color.Cyan("Type your messages (type 'exit' to quit)\n")

	// Initialize message history with system message
	messages = []Message{
		{
			Role:    "system",
			Content: "You are a helpful, accurate, and friendly AI assistant.",
		},
	}

	// Start the prompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("shai> "),
		prompt.OptionTitle("Shai"),
		prompt.OptionInputTextColor(prompt.Cyan),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		// prompt.OptionMultilineCommand([]string{"ctrl+enter", "meta+enter"}),
	)
	p.Run()
}
