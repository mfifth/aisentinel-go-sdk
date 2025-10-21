package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	governor "github.com/aisentinel/aisentinel-go-sdk"
)

func main() {
	cfg := governor.Config{APIKey: "test-key", OfflineMode: true}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	gov, err := governor.NewGovernor(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer gov.Close()

	payload, _ := json.Marshal(map[string]string{"rule-1": "example"})
	result, err := gov.Evaluate(ctx, governor.DecisionRequest{RulepackID: "local", Payload: payload})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("allowed:", result.Allowed, "reason:", result.Reason)
}
