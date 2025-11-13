// 多模型示例程序
// 支持通过命令行参数选择不同的大语言模型：gemini, kimi, qwen, siliconflow, zhipu
// 使用示例：
//
//	go run main.go -model gemini  (使用Gemini模型，默认)
//	go run main.go -model kimi    (使用Kimi模型)
//	go run main.go -model qwen    (使用Qwen模型)
//	go run main.go -model siliconflow  (使用SiliconFlow模型)
//	go run main.go -model zhipu   (使用Zhipu模型)
//	go run main.go -model kimi -model-name kimi-pro  (指定具体的模型名称)
//
// 注意：使用前需设置相应的环境变量来提供API密钥
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/model/kimi"
	"google.golang.org/adk/model/qwen"
	"google.golang.org/adk/model/siliconflow"
	"google.golang.org/adk/model/zhipu"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	// 命令行参数解析
	modelType := flag.String("model", "gemini", "Model type to use: gemini, kimi, qwen, siliconflow, zhipu")
	modelName := flag.String("model-name", "", "Specific model name to use (optional)")
	flag.Parse()

	// 根据模型类型创建相应的模型实例
	var model model.LLM
	var err error

	switch *modelType {
	case "kimi":
		// Kimi模型初始化
		apiKey := os.Getenv("KIMI_API_KEY")
		if apiKey == "" {
			log.Fatal("KIMI_API_KEY environment variable is required")
		}
		config := kimi.Config{APIKey: apiKey}
		modelName := *modelName
		if modelName == "" {
			modelName = kimi.DefaultModel
		}
		model, err = kimi.NewModel(ctx, modelName, config)

	case "qwen":
		// Qwen模型初始化
		apiKey := os.Getenv("QWEN_API_KEY")
		if apiKey == "" {
			log.Fatal("QWEN_API_KEY environment variable is required")
		}
		config := qwen.Config{APIKey: apiKey}
		modelName := *modelName
		if modelName == "" {
			modelName = qwen.DefaultModel
		}
		model, err = qwen.NewModel(ctx, modelName, config)

	case "siliconflow":
		// SiliconFlow模型初始化
		apiKey := os.Getenv("SILICONFLOW_API_KEY")
		if apiKey == "" {
			log.Fatal("SILICONFLOW_API_KEY environment variable is required")
		}
		config := siliconflow.Config{APIKey: apiKey}
		modelName := *modelName
		if modelName == "" {
			modelName = siliconflow.DefaultModel
		}
		model, err = siliconflow.NewModel(ctx, modelName, config)

	case "zhipu":
		// Zhipu模型初始化
		apiKey := os.Getenv("ZHIPU_API_KEY")
		if apiKey == "" {
			log.Fatal("ZHIPU_API_KEY environment variable is required")
		}
		config := zhipu.Config{APIKey: apiKey}
		modelName := *modelName
		if modelName == "" {
			modelName = zhipu.DefaultModel
		}
		model, err = zhipu.NewModel(ctx, modelName, config)

	case "gemini":
		fallthrough
	default:
		// 默认使用Gemini模型
		apiKey := os.Getenv("GOOGLE_API_KEY")
		if apiKey == "" {
			log.Fatal("GOOGLE_API_KEY environment variable is required")
		}
		config := &genai.ClientConfig{APIKey: apiKey}
		modelName := *modelName
		if modelName == "" {
			modelName = "gemini-2.5-flash"
		}
		model, err = gemini.NewModel(ctx, modelName, config)
	}

	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	log.Printf("Successfully initialized %s model\n", *modelType)

	agent, err := llmagent.New(llmagent.Config{
		Name:        "user_assistant",
		Model:       model,
		Description: "Agent to answer user's questions",
		Instruction: "Your SOLE purpose is to answer user's questions",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, flag.Args()); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
