# Shai - LLM Shell

Shai is a command-line shell interface for chatting with Large Language Models using the OpenAI API format. Simply type your queries and get beautifully formatted responses right in your terminal.

## Features

- Interactive shell for conversing with LLMs
- Works with any LLM service that uses the OpenAI API format
- Beautifully renders markdown responses in your terminal
- Maintains conversation history for context-aware responses
- Configurable API endpoints and models

## Installation

### Prerequisites

- Go 1.18 or higher
- Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/alainrk/shai.git
cd shai

# Build and install
make install
```

This will install the `shai` executable to your `$GOPATH/bin` directory, which should be in your PATH.

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/alainrk/shai.git
cd shai

# Build the application
go build -o shai cmd/cli/main.go

# Move to a directory in your PATH
sudo mv shai /usr/local/bin/
```

## Configuration

On first run, Shai creates a default configuration file at `~/.config/shai/config`. Edit this file to add your API key and customize settings:

```
LLM_API_KEY=your_api_key_here

LLM_API_URL=https://api.openai.com/v1/chat/completions
# For Deepseek: https://api.deepseek.com/v1/chat/completions

LLM_MODEL=gpt-4
# For Deepseek: deepseek-chat
```

You can also specify a custom config location using the `--config` flag:

```bash
shai --config /path/to/config
```

## Usage

Start the shell:

```bash
shai
```

Enter your queries at the prompt:

```
shai> What is the capital of the country where Bologna is?
```

Type `exit` or `quit` to exit the shell.

## Dependencies

- [github.com/charmbracelet/glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [github.com/c-bata/go-prompt](https://github.com/c-bata/go-prompt) - Interactive prompt
- [github.com/fatih/color](https://github.com/fatih/color) - Terminal colors
- [github.com/spf13/viper](https://github.com/spf13/viper) - Configuration management

## License

MIT
