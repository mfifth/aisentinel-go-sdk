package governor

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
)

// Rule defines a governance rule compiled for high performance evaluation.
type Rule struct {
	ID          string
	Description string
	Expression  *regexp.Regexp
	Allow       bool
}

// Evaluator performs rule evaluations with concurrency safety.
type Evaluator struct {
	mu    sync.RWMutex
	rules map[string][]Rule
}

// NewEvaluator creates an evaluator instance.
func NewEvaluator() *Evaluator {
	return &Evaluator{rules: make(map[string][]Rule)}
}

// Preload compiles rules for a specific rulepack.
func (e *Evaluator) Preload(rulepackID string, definitions []RuleDefinition) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	rules := make([]Rule, 0, len(definitions))
	for _, def := range definitions {
		re, err := regexp.Compile(def.Pattern)
		if err != nil {
			return fmt.Errorf("compile rule %s: %w", def.ID, err)
		}
		rules = append(rules, Rule{ID: def.ID, Description: def.Description, Expression: re, Allow: def.Allow})
	}
	e.rules[rulepackID] = rules
	return nil
}

// RuleDefinition mirrors rule definitions from rulepacks.
type RuleDefinition struct {
	ID          string
	Description string
	Pattern     string
	Allow       bool
}

// Evaluate evaluates a payload against the provided rulepack.
func (e *Evaluator) Evaluate(ctx context.Context, pack *Rulepack, payload json.RawMessage) (bool, string, error) {
	e.mu.RLock()
	rules, ok := e.rules[pack.ID]
	e.mu.RUnlock()
	if !ok {
		if err := e.Preload(pack.ID, pack.Rules); err != nil {
			return false, "", err
		}
		e.mu.RLock()
		rules = e.rules[pack.ID]
		e.mu.RUnlock()
	}

	var document map[string]any
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &document); err != nil {
			return false, "payload parse error", err
		}
	}

	// Evaluate rules sequentially; this is intentionally simple while enabling
	// future optimisation with goroutines.
	for _, rule := range rules {
		select {
		case <-ctx.Done():
			return false, "context cancelled", ctx.Err()
		default:
		}
		if docValue, ok := document[rule.ID]; ok {
			if str, ok := docValue.(string); ok {
				if rule.Expression.MatchString(str) {
					if rule.Allow {
						return true, rule.Description, nil
					}
					return false, rule.Description, nil
				}
			}
		}
	}

	// Default deny to match Python SDK semantics.
	return false, "no matching rule", nil
}
