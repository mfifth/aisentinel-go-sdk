# AISentinel Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/aisentinel/aisentinel-go-sdk.svg)](https://pkg.go.dev/github.com/aisentinel/aisentinel-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/aisentinel/aisentinel-go-sdk)](https://goreportcard.com/report/github.com/aisentinel/aisentinel-go-sdk)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The official Go SDK for AISentinel, providing comprehensive content moderation and safety evaluation capabilities.

## Features

- **Content Moderation**: Evaluate text, images, and other content for safety and compliance
- **Rule-Based Evaluation**: Flexible rulepack system for customizable moderation policies
- **Caching**: Built-in caching for improved performance and offline capabilities
- **Storage Interfaces**: Pluggable storage backends for reports and evaluation data
- **PII Detection**: Advanced personally identifiable information detection
- **Batch Processing**: Efficient batch evaluation of multiple content items
- **Async Support**: Asynchronous evaluation for high-throughput applications

## Installation

```bash
go get github.com/aisentinel/aisentinel-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aisentinel/aisentinel-go-sdk"
)

func main() {
    // Initialize the client
    client, err := aisentinel.NewClient(aisentinel.Config{
        APIKey: "your-api-key-here",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Evaluate content
    result, err := client.Evaluate(context.Background(), aisentinel.EvaluationRequest{
        Content: "Hello, world!",
        ContentType: aisentinel.ContentTypeText,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Safety Score: %.2f\n", result.Score)
    fmt.Printf("Approved: %t\n", result.Approved)
}
```

## Configuration

```go
config := aisentinel.Config{
    APIKey:      "your-api-key",
    BaseURL:     "https://api.aisentinel.com", // Optional, defaults to production
    Timeout:     30 * time.Second,             // Optional
    MaxRetries:  3,                           // Optional
    CacheConfig: &aisentinel.CacheConfig{     // Optional caching
        Enabled: true,
        TTL:     5 * time.Minute,
    },
}
```

## Usage Examples

### Text Moderation

```go
result, err := client.Evaluate(context.Background(), aisentinel.EvaluationRequest{
    Content:     "User generated text content",
    ContentType: aisentinel.ContentTypeText,
    Rules:       []string{"hate-speech", "toxicity"},
})
```

### Image Moderation

```go
result, err := client.Evaluate(context.Background(), aisentinel.EvaluationRequest{
    Content:     "base64-encoded-image-data",
    ContentType: aisentinel.ContentTypeImage,
    Metadata: map[string]interface{}{
        "filename": "image.jpg",
        "size":     1024000,
    },
})
```

### Batch Evaluation

```go
requests := []aisentinel.EvaluationRequest{
    {Content: "First content", ContentType: aisentinel.ContentTypeText},
    {Content: "Second content", ContentType: aisentinel.ContentTypeText},
}

results, err := client.EvaluateBatch(context.Background(), requests)
```

### Custom Storage Backend

```go
// Implement the Storage interface
type customStorage struct{}

func (s *customStorage) Store(report *aisentinel.Report) error {
    // Custom storage logic
    return nil
}

func (s *customStorage) Retrieve(id string) (*aisentinel.Report, error) {
    // Custom retrieval logic
    return nil, nil
}

// Use with client
client, err := aisentinel.NewClient(aisentinel.Config{
    APIKey:  "your-api-key",
    Storage: &customStorage{},
})
```

## Advanced Features

### Rulepack Management

```go
// Install a custom rulepack
err := client.InstallRulepack(context.Background(), aisentinel.RulepackRequest{
    Name:    "custom-rules",
    Content: rulepackYAML,
})

// List available rulepacks
rulepacks, err := client.ListRulepacks(context.Background())
```

### Offline Mode

```go
// Enable offline evaluation
client.EnableOfflineMode()

// Evaluate without API calls (uses cached rules)
result, err := client.EvaluateOffline(content)
```

### PII Detection

```go
piiResults, err := client.DetectPII(context.Background(), "Text containing email@example.com and phone 555-1234")

for _, detection := range piiResults {
    fmt.Printf("Found %s: %s\n", detection.Type, detection.Value)
}
```

## Error Handling

The SDK provides detailed error information:

```go
result, err := client.Evaluate(context.Background(), request)
if err != nil {
    var apiErr *aisentinel.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (Status: %d)\n", apiErr.Message, apiErr.StatusCode)
    }
    return
}
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://docs.aisentinel.com/go-sdk)
- 🐛 [Issue Tracker](https://github.com/aisentinel/aisentinel-go-sdk/issues)
- 💬 [Community Forum](https://community.aisentinel.com)
- 📧 [Email Support](mailto:support@aisentinel.com)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and updates.
