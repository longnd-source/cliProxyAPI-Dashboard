package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid Config",
			config: Config{
				Port:             8080,
				RequestRetry:     3,
				MaxRetryInterval: 10,
			},
			wantErr: false,
		},
		{
			name: "Invalid Port Low",
			config: Config{
				Port: 0,
			},
			wantErr: true,
			errMsg:  "invalid port: 0",
		},
		{
			name: "Invalid Port High",
			config: Config{
				Port: 70000,
			},
			wantErr: true,
			errMsg:  "invalid port: 70000",
		},
		{
			name: "Invalid RequestRetry",
			config: Config{
				Port:         8080,
				RequestRetry: -1,
			},
			wantErr: true,
			errMsg:  "invalid request-retry: -1",
		},
		{
			name: "Invalid MaxRetryInterval",
			config: Config{
				Port:             8080,
				RequestRetry:     3,
				MaxRetryInterval: -5,
			},
			wantErr: true,
			errMsg:  "invalid max-retry-interval: -5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
