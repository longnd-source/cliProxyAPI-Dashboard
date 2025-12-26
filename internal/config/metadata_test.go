package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestAPIKeyMetadata_YAML(t *testing.T) {
	yamlData := `
api-keys:
  - "key1"
api-key-metadata:
  "key1":
    description: "Production Key"
    tags: ["prod", "critical"]
  "key2":
    description: "Dev Key"
    tags: ["dev"]
`

	var cfg SDKConfig
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	assert.NoError(t, err)

	assert.Contains(t, cfg.APIKeys, "key1")
	assert.Len(t, cfg.APIKeyMetadata, 2)

	meta1, ok := cfg.APIKeyMetadata["key1"]
	assert.True(t, ok)
	assert.Equal(t, "Production Key", meta1.Description)
	assert.Equal(t, []string{"prod", "critical"}, meta1.Tags)

	meta2, ok := cfg.APIKeyMetadata["key2"]
	assert.True(t, ok)
	assert.Equal(t, "Dev Key", meta2.Description)
	assert.Equal(t, []string{"dev"}, meta2.Tags)
}

func TestAPIKeyMetadata_RoundTrip(t *testing.T) {
	cfg := SDKConfig{
		APIKeys: []string{"test-key"},
		APIKeyMetadata: map[string]APIKeyMeta{
			"test-key": {
				Description: "Test Description",
				Tags:        []string{"tag1", "tag2"},
			},
		},
	}

	data, err := yaml.Marshal(&cfg)
	assert.NoError(t, err)

	var loaded SDKConfig
	err = yaml.Unmarshal(data, &loaded)
	assert.NoError(t, err)

	assert.Equal(t, cfg.APIKeys, loaded.APIKeys)
	assert.Equal(t, cfg.APIKeyMetadata, loaded.APIKeyMetadata)
}
