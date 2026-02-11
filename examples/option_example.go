// Option类型使用示例
// 展示如何在Go中使用Rust风格的Option类型进行空安全编程

package main

import (
	"fmt"
	"strings"

	"github.com/dongrv/rust-go"
)

func main() {
	fmt.Println("=== Option类型示例 ===")
	fmt.Println()

	// 1. 基本创建和使用
	fmt.Println("1. 基本创建和使用:")
	exampleBasicOption()
	fmt.Println()

	// 2. 安全解包
	fmt.Println("2. 安全解包:")
	exampleSafeUnwrap()
	fmt.Println()

	// 3. 链式操作
	fmt.Println("3. 链式操作:")
	exampleChaining()
	fmt.Println()

	// 4. 实际应用场景
	fmt.Println("4. 实际应用场景:")
	exampleRealWorld()
	fmt.Println()

	// 5. 组合使用
	fmt.Println("5. 组合使用:")
	exampleCombination()
}

// 1. 基本创建和使用
func exampleBasicOption() {
	// 创建Some值
	someValue := rust.Some(42)
	fmt.Printf("Some(42): %v\n", someValue)
	fmt.Printf("IsSome: %v\n", someValue.IsSome())
	fmt.Printf("IsNone: %v\n", someValue.IsNone())
	fmt.Printf("Unwrap: %v\n", someValue.Unwrap())

	// 创建None值
	noneValue := rust.None[int]()
	fmt.Printf("None[int](): %v\n", noneValue)
	fmt.Printf("IsSome: %v\n", noneValue.IsSome())
	fmt.Printf("IsNone: %v\n", noneValue.IsNone())

	// 使用UnwrapOr提供默认值
	fmt.Printf("Some(42).UnwrapOr(100): %v\n", someValue.UnwrapOr(100))
	fmt.Printf("None[int]().UnwrapOr(100): %v\n", noneValue.UnwrapOr(100))

	// 使用UnwrapOrElse延迟计算默认值
	fmt.Printf("Some(42).UnwrapOrElse(func() int { return 999 }): %v\n",
		someValue.UnwrapOrElse(func() int { return 999 }))
	fmt.Printf("None[int]().UnwrapOrElse(func() int { return 999 }): %v\n",
		noneValue.UnwrapOrElse(func() int { return 999 }))
}

// 2. 安全解包
func exampleSafeUnwrap() {
	// Expect提供自定义panic消息
	safeValue := rust.Some("hello")
	fmt.Printf("Some(\"hello\").Expect(\"should have value\"): %v\n",
		safeValue.Expect("should have value"))

	// 尝试对None使用Expect会panic（在实际代码中应该避免）
	// noneValue := rust.None[string]()
	// 下面的代码会panic: panic: expected string value: no value
	// fmt.Println(noneValue.Expect("expected string value"))

	// 使用Filter过滤值
	evenOption := rust.Some(4).Filter(func(x int) bool { return x%2 == 0 })
	oddOption := rust.Some(3).Filter(func(x int) bool { return x%2 == 0 })
	fmt.Printf("Some(4).Filter(isEven): %v\n", evenOption)
	fmt.Printf("Some(3).Filter(isEven): %v\n", oddOption)

	// Or和OrElse提供备选值
	primary := rust.None[string]()
	fallback := rust.Some("fallback")
	fmt.Printf("None.Or(Some(\"fallback\")): %v\n", primary.Or(fallback))
	fmt.Printf("None.OrElse(func() Option[string] { return Some(\"computed\") }): %v\n",
		primary.OrElse(func() rust.Option[string] { return rust.Some("computed") }))
}

// 3. 链式操作
func exampleChaining() {
	// Map转换值
	number := rust.Some(21)
	doubled := rust.MapOption(number, func(x int) int { return x * 2 })
	fmt.Printf("MapOption(Some(21), x => x*2): %v\n", doubled)

	// AndThen链式调用返回Option的函数
	result := rust.AndThenOption(number, func(x int) rust.Option[string] {
		if x > 10 {
			return rust.Some(fmt.Sprintf("Large number: %d", x))
		}
		return rust.None[string]()
	})
	fmt.Printf("AndThenOption(Some(21), x => Some(\"Large number: ...\")): %v\n", result)

	// 复杂的链式操作
	userInput := rust.Some("  42  ")
	processed := rust.AndThenOption(userInput, func(s string) rust.Option[int] {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return rust.None[int]()
		}
		return rust.Some(len(trimmed))
	})
	fmt.Printf("处理用户输入 \"  42  \" 的结果: %v\n", processed)
}

// 4. 实际应用场景
func exampleRealWorld() {
	// 场景1: 数据库查询可能返回空值
	type User struct {
		ID   int
		Name string
	}

	findUserByID := func(id int) rust.Option[User] {
		// 模拟数据库查询
		if id == 1 {
			return rust.Some(User{ID: 1, Name: "Alice"})
		}
		return rust.None[User]()
	}

	user1 := findUserByID(1)
	user2 := findUserByID(999)

	fmt.Println("查找用户ID 1:", user1)
	fmt.Println("查找用户ID 999:", user2)

	// 安全地使用用户数据
	userName1 := rust.MapOption(user1, func(u User) string { return u.Name })
	userName2 := rust.MapOption(user2, func(u User) string { return u.Name })

	fmt.Printf("用户1的名字: %v\n", userName1.UnwrapOr("未知用户"))
	fmt.Printf("用户2的名字: %v\n", userName2.UnwrapOr("未知用户"))

	// 场景2: 配置解析
	parseConfig := func(config map[string]string) rust.Option[int] {
		if _, ok := config["port"]; ok {
			// 这里可以添加更复杂的解析逻辑
			return rust.Some(8080) // 简化示例
		}
		return rust.None[int]()
	}

	config1 := map[string]string{"port": "8080"}
	config2 := map[string]string{}

	port1 := parseConfig(config1)
	port2 := parseConfig(config2)

	fmt.Printf("配置1的端口: %v\n", port1.UnwrapOr(3000))
	fmt.Printf("配置2的端口: %v\n", port2.UnwrapOr(3000))
}

// 5. 组合使用
func exampleCombination() {
	// 组合多个Option操作
	processNumber := func(input rust.Option[int]) rust.Option[string] {
		return rust.AndThenOption(input, func(n int) rust.Option[string] {
			// 第一步: 检查是否为正数
			if n <= 0 {
				return rust.None[string]()
			}

			// 第二步: 转换为字符串
			str := fmt.Sprintf("Number: %d", n)

			// 第三步: 检查字符串长度
			if len(str) > 10 {
				return rust.None[string]()
			}

			return rust.Some(str)
		})
	}

	testCases := []rust.Option[int]{
		rust.Some(42),
		rust.Some(-5),
		rust.Some(1000000),
		rust.None[int](),
	}

	fmt.Println("组合操作测试:")
	for i, testCase := range testCases {
		result := processNumber(testCase)
		fmt.Printf("  测试%d: 输入=%v, 输出=%v\n", i+1, testCase, result)
	}

	// 使用Option进行错误处理的模式
	fmt.Println("\nOption模式 vs 传统模式:")

	// 传统Go风格
	divideTraditional := func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}

	// RustGo风格
	divideOption := func(a, b int) rust.Option[int] {
		if b == 0 {
			return rust.None[int]()
		}
		return rust.Some(a / b)
	}

	// 比较两种风格
	result1, err1 := divideTraditional(10, 2)
	result2 := divideOption(10, 2)
	result3, err3 := divideTraditional(10, 0)
	result4 := divideOption(10, 0)

	fmt.Printf("  传统风格 10/2: %d, 错误: %v\n", result1, err1)
	fmt.Printf("  Option风格 10/2: %v\n", result2)
	fmt.Printf("  传统风格 10/0: %d, 错误: %v\n", result3, err3)
	fmt.Printf("  Option风格 10/0: %v\n", result4)
}
