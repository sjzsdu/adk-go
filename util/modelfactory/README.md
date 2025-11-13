# Model Factory 模块使用指南

## 概述

Model Factory 模块提供了一个统一的接口，用于创建不同的大语言模型实例，支持多种模型提供商（Gemini、Kimi、Qwen、SiliconFlow、Zhipu）。该模块简化了模型初始化过程，自动处理命令行参数解析和环境变量配置。

## 主要功能

- **统一模型创建接口**：使用一致的API创建不同类型的模型
- **命令行参数支持**：自动识别并处理`-model`和`-model-name`参数
- **环境变量管理**：自动从环境变量获取API密钥
- **默认模型配置**：为每种模型类型提供合理的默认值

## 基本使用方法

### 方法一：使用命令行参数（推荐用于示例程序）

```go
package main

import (
	"context"
	"flag"

	"google.golang.org/adk/util/modelfactory"
)

func main() {
	ctx := context.Background()
	
	// 解析命令行参数（这一步是必须的）
	flag.Parse()
	
	// 使用MustCreateModel自动从命令行参数创建模型
	// 如果创建失败会直接退出程序
	model := modelfactory.MustCreateModel(ctx, nil)
	
	// 使用model进行后续操作...
}
```

### 方法二：手动配置模型参数

```go
package main

import (
	"context"
	
	"google.golang.org/adk/util/modelfactory"
)

func main() {
	ctx := context.Background()
	
	// 手动配置模型参数
	cfg := &modelfactory.Config{
		ModelType: "kimi",     // 指定模型类型
		ModelName: "kimi-pro", // 指定具体模型名称（可选）
	}
	
	// 创建模型，带错误处理
	model, err := modelfactory.CreateModel(ctx, cfg)
	if err != nil {
		// 处理错误
		panic(err)
	}
	
	// 使用model进行后续操作...
}
```

### 方法三：自定义命令行参数后使用工厂

```go
package main

import (
	"context"
	"flag"
	
	"google.golang.org/adk/util/modelfactory"
)

func main() {
	ctx := context.Background()
	
	// 定义自定义标志（如果需要的话）
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	
	// 让modelfactory处理模型相关的标志
	cfg := modelfactory.NewFromFlags()
	
	// 解析所有标志
	flag.Parse()
	
	// 创建模型
	model := modelfactory.MustCreateModel(ctx, cfg)
	
	// 使用model进行后续操作...
}
```

## 在现有示例中集成

要在现有的示例程序中使用Model Factory，请按照以下步骤操作：

1. **导入模块**：添加对`modelfactory`包的导入
2. **移除原有模型初始化代码**：删除重复的模型创建逻辑
3. **使用工厂方法**：调用`MustCreateModel`或`CreateModel`方法
4. **保留flag.Parse()**：确保在使用工厂前调用`flag.Parse()`

## 环境变量要求

使用不同模型时，需要设置相应的环境变量：

- **Gemini**: `GOOGLE_API_KEY`
- **Kimi**: `KIMI_API_KEY`
- **Qwen**: `QWEN_API_KEY`
- **SiliconFlow**: `SILICONFLOW_API_KEY`
- **Zhipu**: `ZHIPU_API_KEY`

## 命令行参数

Model Factory支持以下命令行参数：

- `-model`: 指定要使用的模型类型（gemini, kimi, qwen, siliconflow, zhipu），默认为gemini
- `-model-name`: 指定具体的模型名称（可选），如果不指定则使用默认模型

## 示例

### 在quickstart示例中使用

```go
package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/util/modelfactory"
)

func main() {
	ctx := context.Background()
	flag.Parse()
	
	// 使用Model Factory创建模型
	model := modelfactory.MustCreateModel(ctx, nil)
	
	// 创建代理
	agent, err := llmagent.New(llmagent.Config{
		Name:        "quickstart_agent",
		Model:       model,
		Description: "Quickstart agent",
		Instruction: "You are a helpful assistant.",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}
	
	// 启动服务
	config := &launcher.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, flag.Args()); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
```

## 注意事项

1. **标志定义冲突**：如果您的程序已经定义了`-model`或`-model-name`标志，Model Factory会自动使用这些已定义的标志，而不会重新定义它们。

2. **错误处理**：`MustCreateModel`方法在创建失败时会直接调用`log.Fatal`，适合示例程序；而`CreateModel`方法会返回错误，适合需要自定义错误处理的场景。

3. **环境变量**：确保在运行程序前设置了正确的环境变量，否则模型创建会失败。