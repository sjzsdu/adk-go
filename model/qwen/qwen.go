package qwen

// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package qwen implements the [model.LLM] interface for Qwen (通义千问) models using OpenAI-compatible API.

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/openai"
)

const (
	// 环境变量名
	TokenEnvVarName = "QWEN_API_KEY" //nolint:gosec
	ModelEnvVarName = "QWEN_MODEL"   //nolint:gosec

	// OpenAI兼容模式基础URL
	DefaultBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"

	// 默认模型
	DefaultModel = "qwen-max"
)

// Qwen模型常量
const (
	// ModelQWenTurbo 是通义千问Turbo模型
	ModelQWenTurbo = "qwen-turbo"

	// ModelQWenPlus 是通义千问Plus模型
	ModelQWenPlus = "qwen-plus"

	// ModelQWenMax 是通义千问Max模型
	ModelQWenMax = "qwen-max"

	// ModelQWenVLPlus 是通义千问视觉Plus模型
	ModelQWenVLPlus = "qwen-vl-plus"

	// ModelQWenVLMax 是通义千问视觉Max模型
	ModelQWenVLMax = "qwen-vl-max"
)

// Config holds the configuration for Qwen model initialization.
type Config struct {
	// APIKey is the Qwen API key. If empty, it will be read from QWEN_API_KEY environment variable.
	APIKey string
	// BaseURL is the Qwen API base URL. If empty, it will use DefaultBaseURL.
	BaseURL string
	// Organization is the organization ID (optional for Qwen).
	Organization string
}

// NewModel returns [model.LLM], backed by the Qwen API using OpenAI-compatible interface.
//
// It uses the provided modelName and configuration to initialize the underlying
// OpenAI client with Qwen-specific endpoints and authentication.
//
// An error is returned if the configuration is invalid.
func NewModel(ctx context.Context, modelName string, config Config) (model.LLM, error) {
	// Set defaults
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(TokenEnvVarName)
		if apiKey == "" {
			return nil, fmt.Errorf("Qwen API key is required, set QWEN_API_KEY environment variable or provide APIKey in config")
		}
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Validate model name
	if !isValidQwenModel(modelName) {
		return nil, fmt.Errorf("unsupported Qwen model: %s", modelName)
	}

	// Create OpenAI config with Qwen-specific settings
	openaiConfig := openai.Config{
		APIKey:       apiKey,
		BaseURL:      baseURL,
		Organization: config.Organization,
	}

	// Use OpenAI implementation with Qwen configuration
	return openai.NewModel(ctx, modelName, openaiConfig)
}

// isValidQwenModel checks if the given model name is a valid Qwen model.
func isValidQwenModel(modelName string) bool {
	validModels := map[string]bool{
		ModelQWenTurbo:  true,
		ModelQWenPlus:   true,
		ModelQWenMax:    true,
		ModelQWenVLPlus: true,
		ModelQWenVLMax:  true,
	}
	return validModels[modelName]
}

// GetSupportedModels returns a list of supported Qwen models.
func GetSupportedModels() []string {
	return []string{
		ModelQWenTurbo,
		ModelQWenPlus,
		ModelQWenMax,
		ModelQWenVLPlus,
		ModelQWenVLMax,
	}
}

// getEnvOrDefault 获取环境变量值，如果不存在则返回默认值
func getEnvOrDefault(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}
	return value
}
