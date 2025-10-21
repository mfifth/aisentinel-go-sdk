package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	governor "github.com/mfifth/aisentinel-go-sdk"
	"github.com/mfifth/aisentinel-go-sdk/storage"
)

func main() {
	cfg := governor.Config{APIKey: "demo", StorageBackend: string(storage.BackendMemory), CacheTTL: time.Minute}
	ctx := context.Background()

	gov, err := governor.NewGovernor(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer gov.Close()

	payload, _ := json.Marshal(map[string]any{"rule-1": "sensitive"})
	result, err := gov.Evaluate(ctx, governor.DecisionRequest{RulepackID: "governance", Payload: payload})
	if err != nil {
		fmt.Println("decision error:", err)
		return
	}
	fmt.Printf("allowed=%v reason=%s latency=%s\n", result.Allowed, result.Reason, result.Latency)
}
