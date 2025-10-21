package governor

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config encapsulates runtime configuration for the Governor. It mirrors the
// Python SDK configuration surface while staying idiomatic to Go.
type Config struct {
	APIBaseURL        string
	APIKey            string
	CacheTTL          time.Duration
	HTTPTimeout       time.Duration
	OfflineMode       bool
	OfflineQueueSize  int
	StorageBackend    string
	StorageDSN        string
	MetricsEnabled    bool
	MetricsEndpoint   string
	EnvironmentPrefix string
}

// DefaultConfig returns a configuration populated with production ready defaults.
func DefaultConfig() Config {
	return Config{
		APIBaseURL:        "https://api.aisentinel.ai",
		CacheTTL:          5 * time.Minute,
		HTTPTimeout:       10 * time.Second,
		OfflineMode:       false,
		OfflineQueueSize:  1024,
		StorageBackend:    "memory",
		MetricsEnabled:    true,
		EnvironmentPrefix: "AISENTINEL_",
	}
}

// ApplyEnv overlays configuration values from environment variables using the
// configured prefix. The behaviour matches the Python SDK to ease migration.
func (c *Config) ApplyEnv() error {
	prefix := c.EnvironmentPrefix
	if prefix == "" {
		prefix = "AISENTINEL_"
	}
	overlay := map[string]func(string) error{
		"API_BASE_URL": func(v string) error {
			if _, err := url.ParseRequestURI(v); err != nil {
				return fmt.Errorf("invalid API_BASE_URL: %w", err)
			}
			c.APIBaseURL = v
			return nil
		},
		"API_KEY": func(v string) error {
			c.APIKey = v
			return nil
		},
		"CACHE_TTL": func(v string) error {
			d, err := time.ParseDuration(v)
			if err != nil {
				return fmt.Errorf("invalid CACHE_TTL: %w", err)
			}
			c.CacheTTL = d
			return nil
		},
		"HTTP_TIMEOUT": func(v string) error {
			d, err := time.ParseDuration(v)
			if err != nil {
				return fmt.Errorf("invalid HTTP_TIMEOUT: %w", err)
			}
			c.HTTPTimeout = d
			return nil
		},
		"OFFLINE_MODE": func(v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("invalid OFFLINE_MODE: %w", err)
			}
			c.OfflineMode = b
			return nil
		},
		"OFFLINE_QUEUE_SIZE": func(v string) error {
			i, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("invalid OFFLINE_QUEUE_SIZE: %w", err)
			}
			if i <= 0 {
				return fmt.Errorf("offline queue size must be > 0")
			}
			c.OfflineQueueSize = i
			return nil
		},
		"STORAGE_BACKEND": func(v string) error {
			c.StorageBackend = strings.ToLower(v)
			return nil
		},
		"STORAGE_DSN": func(v string) error {
			c.StorageDSN = v
			return nil
		},
		"METRICS_ENABLED": func(v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("invalid METRICS_ENABLED: %w", err)
			}
			c.MetricsEnabled = b
			return nil
		},
		"METRICS_ENDPOINT": func(v string) error {
			c.MetricsEndpoint = v
			return nil
		},
	}

	for key, fn := range overlay {
		if value, ok := os.LookupEnv(prefix + key); ok {
			if err := fn(value); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate performs sanity checks on the configuration.
func (c Config) Validate() error {
	if c.APIBaseURL == "" {
		return fmt.Errorf("APIBaseURL is required")
	}
	if _, err := url.ParseRequestURI(c.APIBaseURL); err != nil {
		return fmt.Errorf("invalid APIBaseURL: %w", err)
	}
	if c.APIKey == "" {
		return fmt.Errorf("APIKey is required")
	}
	if c.CacheTTL <= 0 {
		return fmt.Errorf("CacheTTL must be > 0")
	}
	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTPTimeout must be > 0")
	}
	if c.OfflineQueueSize <= 0 {
		return fmt.Errorf("OfflineQueueSize must be > 0")
	}
	return nil
}

// Merge merges another config into the current one, overriding non-zero values.
func (c Config) Merge(other Config) Config {
	if other.APIBaseURL != "" {
		c.APIBaseURL = other.APIBaseURL
	}
	if other.APIKey != "" {
		c.APIKey = other.APIKey
	}
	if other.CacheTTL != 0 {
		c.CacheTTL = other.CacheTTL
	}
	if other.HTTPTimeout != 0 {
		c.HTTPTimeout = other.HTTPTimeout
	}
	if other.OfflineQueueSize != 0 {
		c.OfflineQueueSize = other.OfflineQueueSize
	}
	if other.StorageBackend != "" {
		c.StorageBackend = other.StorageBackend
	}
	if other.StorageDSN != "" {
		c.StorageDSN = other.StorageDSN
	}
	if other.MetricsEndpoint != "" {
		c.MetricsEndpoint = other.MetricsEndpoint
	}
	if other.EnvironmentPrefix != "" {
		c.EnvironmentPrefix = other.EnvironmentPrefix
	}
	c.OfflineMode = other.OfflineMode
	c.MetricsEnabled = other.MetricsEnabled
	return c
}
