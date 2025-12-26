package config

import "fmt"

func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d; must be between 1 and 65535", c.Port)
	}
	if c.RequestRetry < 0 {
		return fmt.Errorf("invalid request-retry: %d; must be non-negative", c.RequestRetry)
	}
	if c.MaxRetryInterval < 0 {
		return fmt.Errorf("invalid max-retry-interval: %d; must be non-negative", c.MaxRetryInterval)
	}
	return nil
}
