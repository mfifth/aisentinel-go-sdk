package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	aisentinel "github.com/mfifth/aisentinel-go-sdk"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = ""
)

func main() {
	apiKey := flag.String("api-key", os.Getenv("AISENTINEL_API_KEY"), "AISentinel API key (or set AISENTINEL_API_KEY)")
	apiBaseURL := flag.String("api-base-url", "", "Override the AISentinel API base URL")
	rulepack := flag.String("rulepack", "default", "Rulepack identifier to evaluate")
	payloadInline := flag.String("payload", "", "Inline JSON payload to evaluate")
	payloadFile := flag.String("payload-file", "", "Path to a file containing JSON payload")
	offline := flag.Bool("offline", false, "Enable offline evaluation mode")
	timeout := flag.Duration("timeout", 15*time.Second, "Timeout for the evaluation request")
	showVersion := flag.Bool("version", false, "Print version information and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("aisentinel-go-sdk %s (commit %s, built %s)\n", buildVersion, buildCommit, buildDate)
		return
	}

	if *payloadInline != "" && *payloadFile != "" {
		log.Fatal("only one of --payload or --payload-file may be provided")
	}

	if *apiKey == "" {
		log.Fatal("API key is required (set --api-key or AISENTINEL_API_KEY)")
	}

	payload, err := resolvePayload(*payloadInline, *payloadFile)
	if err != nil {
		log.Fatalf("resolve payload: %v", err)
	}

	cfg := aisentinel.Config{ // nolint:exhaustruct
		APIKey:      *apiKey,
		OfflineMode: *offline,
	}
	if *apiBaseURL != "" {
		cfg.APIBaseURL = *apiBaseURL
	}
	if *timeout > 0 {
		cfg.HTTPTimeout = *timeout
	}

	ctx := context.Background()

	governor, err := aisentinel.NewGovernor(ctx, cfg)
	if err != nil {
		log.Fatalf("initialise governor: %v", err)
	}
	defer func() {
		if cerr := governor.Close(); cerr != nil {
			log.Printf("close governor: %v", cerr)
		}
	}()

	evalCtx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()

	result, err := governor.Evaluate(evalCtx, aisentinel.DecisionRequest{ // nolint:exhaustruct
		RulepackID: *rulepack,
		Payload:    payload,
	})
	if err != nil {
		log.Fatalf("evaluate: %v", err)
	}

	output := map[string]any{
		"allowed":    result.Allowed,
		"reason":     result.Reason,
		"latency_ms": result.Latency.Milliseconds(),
	}

	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("encode result: %v", err)
	}
	fmt.Println(string(encoded))
}

func resolvePayload(inline, path string) (json.RawMessage, error) {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if !json.Valid(data) {
			return nil, errors.New("payload file does not contain valid JSON")
		}
		return json.RawMessage(data), nil
	}

	if inline == "" {
		if flag.NArg() > 0 {
			inline = flag.Arg(0)
		} else {
			inline = "{}"
		}
	}

	data := []byte(inline)
	if !json.Valid(data) {
		return nil, errors.New("payload must be valid JSON")
	}
	return json.RawMessage(data), nil
}
