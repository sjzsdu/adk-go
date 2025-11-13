package modelfactory

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/model/kimi"
	"google.golang.org/adk/model/qwen"
	"google.golang.org/adk/model/siliconflow"
	"google.golang.org/adk/model/zhipu"
	"google.golang.org/genai"
)

// 定义全局变量来存储标志值
var (
	// modelType 存储命令行中的模型类型
	modelType = flag.String("model", "gemini", "Model type to use: gemini, kimi, qwen, siliconflow, zhipu")
	// modelName 存储命令行中的模型名称
	modelName = flag.String("model-name", "", "Specific model name to use (optional)")
)

// init 函数在包初始化时自动注册标志
func init() {
	// 标志已经通过全局变量定义，这里可以添加其他初始化逻辑
}

// Config contains model factory configuration options
// that can be set programmatically or via command line flags.
type Config struct {
	ModelType string // Model type to use: gemini, kimi, qwen, siliconflow, zhipu
	ModelName string // Specific model name to use (optional)
}

// NewFromFlags creates a new Config from command line flags.
// Now it simply returns the current values of the globally registered flags.
func NewFromFlags() *Config {
	return &Config{
		ModelType: *modelType,
		ModelName: *modelName,
	}
}

// CreateModel creates a new LLM model based on the provided configuration.
func CreateModel(ctx context.Context, cfg *Config) (model.LLM, error) {
	// 如果没有提供配置，使用从命令行标志获取的值
	if cfg == nil {
		cfg = NewFromFlags()
	}

	log.Printf("Creating %s model...", cfg.ModelType)

	var model model.LLM
	var err error

	// 使用配置中的模型名称或回退到默认值
	modelName := cfg.ModelName

	switch cfg.ModelType {
	case "kimi":
		// Kimi model initialization
		apiKey := os.Getenv("KIMI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("KIMI_API_KEY environment variable is required")
		}
		config := kimi.Config{APIKey: apiKey}
		if modelName == "" {
			modelName = kimi.DefaultModel
		}
		model, err = kimi.NewModel(ctx, modelName, config)

	case "qwen":
		// Qwen model initialization
		apiKey := os.Getenv("QWEN_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("QWEN_API_KEY environment variable is required")
		}
		config := qwen.Config{APIKey: apiKey}
		if modelName == "" {
			modelName = qwen.DefaultModel
		}
		model, err = qwen.NewModel(ctx, modelName, config)

	case "siliconflow":
		// SiliconFlow model initialization
		apiKey := os.Getenv("SILICONFLOW_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("SILICONFLOW_API_KEY environment variable is required")
		}
		config := siliconflow.Config{APIKey: apiKey}
		if modelName == "" {
			modelName = siliconflow.DefaultModel
		}
		model, err = siliconflow.NewModel(ctx, modelName, config)

	case "zhipu":
		// Zhipu model initialization
		apiKey := os.Getenv("ZHIPU_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ZHIPU_API_KEY environment variable is required")
		}
		config := zhipu.Config{APIKey: apiKey}
		if modelName == "" {
			modelName = zhipu.DefaultModel
		}
		model, err = zhipu.NewModel(ctx, modelName, config)

	case "gemini":
		fallthrough
	default:
		// Default to Gemini model
		apiKey := os.Getenv("GOOGLE_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GOOGLE_API_KEY environment variable is required")
		}
		config := &genai.ClientConfig{APIKey: apiKey}
		if modelName == "" {
			modelName = "gemini-2.5-flash"
		}
		model, err = gemini.NewModel(ctx, modelName, config)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s model: %w", cfg.ModelType, err)
	}

	log.Printf("Successfully initialized %s model (name: %s)", cfg.ModelType, modelName)
	return model, nil
}

// MustCreateModel creates a new LLM model and panics on error.
// Useful for examples and quickstart applications where error handling is simplified.
func MustCreateModel(ctx context.Context, cfg *Config) model.LLM {
	model, err := CreateModel(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}
	return model
}
