// Result类型使用示例
// 展示如何在Go中使用Rust风格的Result类型进行错误处理

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dongrv/rust-go"
)

func main() {
	fmt.Println("=== Result类型示例 ===")
	fmt.Println()

	// 1. 基本创建和使用
	fmt.Println("1. 基本创建和使用:")
	exampleBasicResult()
	fmt.Println()

	// 2. 错误处理模式
	fmt.Println("2. 错误处理模式:")
	exampleErrorHandling()
	fmt.Println()

	// 3. 链式操作（铁路编程）
	fmt.Println("3. 链式操作（铁路编程）:")
	exampleRailwayProgramming()
	fmt.Println()

	// 4. 实际应用场景
	fmt.Println("4. 实际应用场景:")
	exampleRealWorld()
	fmt.Println()

	// 5. 高级用法
	fmt.Println("5. 高级用法:")
	exampleAdvancedUsage()
}

// 1. 基本创建和使用
func exampleBasicResult() {
	// 创建Ok值
	okResult := rust.Ok[int, string](42)
	fmt.Printf("Ok[int, string](42): %v\n", okResult)
	fmt.Printf("IsOk: %v\n", okResult.IsOk())
	fmt.Printf("IsErr: %v\n", okResult.IsErr())
	fmt.Printf("Unwrap: %v\n", okResult.Unwrap())

	// 创建Err值
	errResult := rust.Err[int, string]("something went wrong")
	fmt.Printf("Err[int, string](\"something went wrong\"): %v\n", errResult)
	fmt.Printf("IsOk: %v\n", errResult.IsOk())
	fmt.Printf("IsErr: %v\n", errResult.IsErr())
	fmt.Printf("UnwrapErr: %v\n", errResult.UnwrapErr())

	// 安全解包
	fmt.Printf("Ok(42).UnwrapOr(100): %v\n", okResult.UnwrapOr(100))
	fmt.Printf("Err(\"error\").UnwrapOr(100): %v\n", errResult.UnwrapOr(100))

	// 使用UnwrapOrElse处理错误
	fmt.Printf("Ok(42).UnwrapOrElse(func(e string) int { return len(e) }): %v\n",
		okResult.UnwrapOrElse(func(e string) int { return len(e) }))
	fmt.Printf("Err(\"error\").UnwrapOrElse(func(e string) int { return len(e) }): %v\n",
		errResult.UnwrapOrElse(func(e string) int { return len(e) }))
}

// 2. 错误处理模式
func exampleErrorHandling() {
	// Expect提供自定义panic消息
	success := rust.Ok[string, string]("data loaded")
	fmt.Printf("Ok(\"data loaded\").Expect(\"should succeed\"): %v\n",
		success.Expect("should succeed"))

	// 错误恢复模式
	parseInt := func(s string) rust.Result[int, string] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return rust.Err[int, string](fmt.Sprintf("解析失败: %v", err))
		}
		return rust.Ok[int, string](n)
	}

	testCases := []string{"42", "not-a-number", "-100", "3.14"}

	fmt.Println("解析整数测试:")
	for _, test := range testCases {
		result := parseInt(test)
		fmt.Printf("  解析 %q: %v\n", test, result)

		// 使用Or提供备选结果
		fallback := rust.Ok[int, string](0)
		withFallback := result.Or(fallback)
		fmt.Printf("    使用备选值: %v\n", withFallback)

		// 使用OrElse动态生成备选
		dynamicFallback := result.OrElse(func(err string) rust.Result[int, string] {
			return rust.Ok[int, string](-999)
		})
		fmt.Printf("    动态备选值: %v\n", dynamicFallback)
	}

	// MapErr转换错误类型
	originalErr := rust.Err[int, string]("原始错误")
	wrappedErr := rust.MapErrResult(originalErr, func(e string) string {
		return fmt.Sprintf("包装后的错误: %s", e)
	})
	fmt.Printf("\n错误包装: %v -> %v\n", originalErr, wrappedErr)
}

// 3. 链式操作（铁路编程）
func exampleRailwayProgramming() {
	// 定义一些可能失败的操作
	validateInput := func(input string) rust.Result[string, string] {
		if len(input) == 0 {
			return rust.Err[string, string]("输入不能为空")
		}
		if len(input) > 100 {
			return rust.Err[string, string]("输入过长")
		}
		return rust.Ok[string, string](input)
	}

	parseNumber := func(input string) rust.Result[int, string] {
		n, err := strconv.Atoi(input)
		if err != nil {
			return rust.Err[int, string](fmt.Sprintf("无效的数字: %s", input))
		}
		return rust.Ok[int, string](n)
	}

	validateRange := func(n int) rust.Result[int, string] {
		if n < 0 {
			return rust.Err[int, string]("数字不能为负数")
		}
		if n > 1000 {
			return rust.Err[int, string]("数字太大")
		}
		return rust.Ok[int, string](n)
	}

	processNumber := func(n int) rust.Result[string, string] {
		return rust.Ok[string, string](fmt.Sprintf("处理后的数字: %d", n*2))
	}

	// 铁路编程：一系列操作，任何一个失败都会短路
	fmt.Println("铁路编程示例:")

	pipeline := func(input string) rust.Result[string, string] {
		// 第一步：验证输入
		inputResult := validateInput(input)
		if inputResult.IsErr() {
			return rust.Err[string, string](inputResult.UnwrapErr())
		}

		// 第二步：解析数字
		numberResult := parseNumber(inputResult.Unwrap())
		if numberResult.IsErr() {
			return rust.Err[string, string](numberResult.UnwrapErr())
		}

		// 第三步：验证范围
		rangeResult := validateRange(numberResult.Unwrap())
		if rangeResult.IsErr() {
			return rust.Err[string, string](rangeResult.UnwrapErr())
		}

		// 第四步：处理数字
		return processNumber(rangeResult.Unwrap())
	}

	testInputs := []string{"", "abc", "-50", "5000", "42", "100"}

	for _, input := range testInputs {
		result := pipeline(input)
		fmt.Printf("  输入 %q: %v\n", input, result)
	}

	// 使用链式操作简化
	fmt.Println("\n使用链式操作简化:")

	simplePipeline := func(input string) rust.Result[string, string] {
		return rust.MapResult(
			rust.AndThenResult(
				rust.AndThenResult(
					validateInput(input),
					parseNumber,
				),
				validateRange,
			),
			func(n int) string {
				return fmt.Sprintf("简化结果: %d", n*3)
			},
		)
	}

	for _, input := range []string{"25", "invalid"} {
		result := simplePipeline(input)
		fmt.Printf("  输入 %q: %v\n", input, result)
	}
}

// 4. 实际应用场景
func exampleRealWorld() {
	// 场景1: 文件操作
	type FileContent struct {
		Path    string
		Content string
	}

	readFile := func(path string) rust.Result[FileContent, string] {
		// 模拟文件读取
		if path == "" {
			return rust.Err[FileContent, string]("文件路径为空")
		}
		if !strings.HasSuffix(path, ".txt") {
			return rust.Err[FileContent, string]("只支持.txt文件")
		}

		// 模拟成功读取
		return rust.Ok[FileContent, string](FileContent{
			Path:    path,
			Content: fmt.Sprintf("这是 %s 的内容", path),
		})
	}

	processContent := func(content FileContent) rust.Result[string, string] {
		if len(content.Content) > 1000 {
			return rust.Err[string, string]("文件内容过长")
		}
		return rust.Ok[string, string](strings.ToUpper(content.Content))
	}

	fmt.Println("文件处理流程:")

	fileProcessing := func(path string) rust.Result[string, string] {
		return rust.AndThenResult(
			readFile(path),
			processContent,
		)
	}

	fileTests := []string{"", "data.pdf", "document.txt", "large.txt"}

	for _, file := range fileTests {
		result := fileProcessing(file)
		fmt.Printf("  处理文件 %q: %v\n", file, result)

		// 优雅地处理结果
		finalOutput := result.UnwrapOrElse(func(err string) string {
			return fmt.Sprintf("错误: %s", err)
		})
		fmt.Printf("    最终输出: %s\n", finalOutput)
	}

	// 场景2: API调用
	fmt.Println("\nAPI调用示例:")

	type APIResponse struct {
		Status  int
		Data    string
		Message string
	}

	callAPI := func(endpoint string) rust.Result[APIResponse, string] {
		// 模拟API调用
		switch endpoint {
		case "/users":
			return rust.Ok[APIResponse, string](APIResponse{
				Status:  200,
				Data:    `[{"id": 1, "name": "Alice"}]`,
				Message: "成功",
			})
		case "/products":
			return rust.Ok[APIResponse, string](APIResponse{
				Status:  200,
				Data:    `[{"id": 1, "name": "Product A"}]`,
				Message: "成功",
			})
		default:
			return rust.Err[APIResponse, string]("端点不存在")
		}
	}

	parseResponse := func(resp APIResponse) rust.Result[map[string]interface{}, string] {
		if resp.Status != 200 {
			return rust.Err[map[string]interface{}, string](
				fmt.Sprintf("API返回错误状态: %d", resp.Status),
			)
		}
		// 这里可以添加JSON解析逻辑
		return rust.Ok[map[string]interface{}, string](map[string]interface{}{
			"data":    resp.Data,
			"message": resp.Message,
		})
	}

	apiPipeline := func(endpoint string) rust.Result[map[string]interface{}, string] {
		return rust.AndThenResult(
			callAPI(endpoint),
			parseResponse,
		)
	}

	endpoints := []string{"/users", "/products", "/orders"}

	for _, endpoint := range endpoints {
		result := apiPipeline(endpoint)
		fmt.Printf("  调用 %s: %v\n", endpoint, result)
	}
}

// 5. 高级用法
func exampleAdvancedUsage() {
	// 组合Result和Option
	fmt.Println("组合Result和Option:")

	// 一个返回Result的函数
	fetchData := func(id int) rust.Result[string, string] {
		if id < 0 {
			return rust.Err[string, string]("无效的ID")
		}
		if id > 100 {
			return rust.Ok[string, string]("")
		}
		return rust.Ok[string, string](fmt.Sprintf("数据-%d", id))
	}

	// 一个返回Option的函数
	extractValue := func(data string) rust.Option[int] {
		if data == "" {
			return rust.None[int]()
		}
		// 简单示例：返回数据长度
		return rust.Some(len(data))
	}

	// 组合使用
	combinedPipeline := func(id int) rust.Option[int] {
		dataResult := fetchData(id)
		if dataResult.IsErr() {
			return rust.None[int]()
		}
		return extractValue(dataResult.Unwrap())
	}

	// 更优雅的方式
	betterPipeline := func(id int) rust.Option[int] {
		dataResult := fetchData(id)
		if dataResult.IsErr() {
			return rust.None[int]()
		}
		return extractValue(dataResult.Unwrap())
	}

	testIDs := []int{-1, 0, 50, 150}

	fmt.Println("组合管道测试:")
	for _, id := range testIDs {
		result1 := combinedPipeline(id)
		result2 := betterPipeline(id)
		fmt.Printf("  ID %d: 方法1=%v, 方法2=%v\n", id, result1, result2)
	}

	// 错误处理策略比较
	fmt.Println("\n错误处理策略比较:")

	// 传统Go风格
	processTraditional := func() (string, error) {
		// 多个可能失败的操作
		if err := step1(); err != nil {
			return "", fmt.Errorf("步骤1失败: %w", err)
		}

		data, err := step2()
		if err != nil {
			return "", fmt.Errorf("步骤2失败: %w", err)
		}

		result, err := step3(data)
		if err != nil {
			return "", fmt.Errorf("步骤3失败: %w", err)
		}

		return result, nil
	}

	// RustGo风格
	processRustGo := func() rust.Result[string, string] {
		return rust.AndThenResult(
			rust.AndThenResult(
				step1Result(),
				func(_ struct{}) rust.Result[string, string] {
					return step2Result()
				},
			),
			step3Result,
		)
	}

	// 模拟函数
	_ = processTraditional
	_ = processRustGo

	fmt.Println("  传统风格: 需要多次错误检查")
	fmt.Println("  RustGo风格: 链式操作，错误自动传播")
	fmt.Println("  优点: 更清晰的错误流，减少嵌套if语句")
}

// 模拟函数实现
func step1() error                      { return nil }
func step2() (string, error)            { return "data", nil }
func step3(data string) (string, error) { return "result", nil }

func step1Result() rust.Result[struct{}, string] {
	return rust.Ok[struct{}, string](struct{}{})
}

func step2Result() rust.Result[string, string] {
	return rust.Ok[string, string]("data")
}

func step3Result(data string) rust.Result[string, string] {
	return rust.Ok[string, string]("result")
}
