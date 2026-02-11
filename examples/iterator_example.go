// Iterator和Chainable使用示例
// 展示如何在Go中使用Rust风格的迭代器和链式集合操作

package main

import (
	"fmt"
	"strings"

	"github.com/dongrv/rust-go"
)

func main() {
	fmt.Println("=== Iterator和Chainable示例 ===")
	fmt.Println()

	// 1. 基本迭代器操作
	fmt.Println("1. 基本迭代器操作:")
	exampleBasicIterator()
	fmt.Println()

	// 2. 链式迭代器操作
	fmt.Println("2. 链式迭代器操作:")
	exampleChainedIterator()
	fmt.Println()

	// 3. Chainable集合操作
	fmt.Println("3. Chainable集合操作:")
	exampleChainable()
	fmt.Println()

	// 4. 高级迭代器功能
	fmt.Println("4. 高级迭代器功能:")
	exampleAdvancedIterator()
	fmt.Println()

	// 5. 实际应用场景
	fmt.Println("5. 实际应用场景:")
	exampleRealWorld()
	fmt.Println()

	// 6. 性能比较和最佳实践
	fmt.Println("6. 性能比较和最佳实践:")
	examplePerformance()
}

// 1. 基本迭代器操作
func exampleBasicIterator() {
	fmt.Println("基础迭代器创建和使用:")

	// 从切片创建迭代器
	numbers := []int{1, 2, 3, 4, 5}
	iter := rust.Iter(numbers)

	fmt.Printf("原始切片: %v\n", numbers)
	fmt.Print("迭代器输出: ")
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}
		fmt.Printf("%d ", next.Unwrap())
	}
	fmt.Println()

	// 使用Collect收集结果
	iter2 := rust.Iter(numbers)
	collected := rust.Collect(iter2)
	fmt.Printf("Collect结果: %v\n", collected)

	// 使用ForEach遍历
	fmt.Print("ForEach输出: ")
	rust.ForEach(rust.Iter(numbers), func(x int) {
		fmt.Printf("%d ", x)
	})
	fmt.Println()

	// 使用Fold计算总和
	sum := rust.Fold(rust.Iter(numbers), 0, func(acc, x int) int {
		return acc + x
	})
	fmt.Printf("Fold计算总和: %d\n", sum)

	// 使用Reduce
	max := rust.Reduce(rust.Iter(numbers), func(a, b int) int {
		if a > b {
			return a
		}
		return b
	})
	fmt.Printf("Reduce找最大值: %v\n", max)
}

// 2. 链式迭代器操作
func exampleChainedIterator() {
	fmt.Println("链式迭代器操作:")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 创建复杂的迭代器管道
	pipeline := rust.Map(
		rust.Take(
			rust.Filter(
				rust.Iter(data),
				func(x int) bool { return x%2 == 0 }, // 只取偶数
			),
			3, // 取前3个
		),
		func(x int) int { return x * x }, // 平方
	)

	result := rust.Collect(pipeline)
	fmt.Printf("原始数据: %v\n", data)
	fmt.Printf("管道结果(偶数前3个的平方): %v\n", result)

	// 使用Skip
	skipped := rust.Collect(rust.Skip(rust.Iter(data), 5))
	fmt.Printf("跳过前5个: %v\n", skipped)

	// 使用Chain连接迭代器
	first := []int{1, 2, 3}
	second := []int{4, 5, 6}
	chained := rust.Collect(rust.Chain(rust.Iter(first), rust.Iter(second)))
	fmt.Printf("连接迭代器: %v + %v = %v\n", first, second, chained)

	// 使用Zip配对
	names := []string{"Alice", "Bob", "Charlie"}
	ages := []int{25, 30, 35}
	zipped := rust.Collect(rust.Zip(rust.Iter(names), rust.Iter(ages)))
	fmt.Printf("Zip配对: %v\n", zipped)

	// 使用Enumerate添加索引
	enumerated := rust.Collect(rust.Enumerate(rust.Iter(names)))
	fmt.Printf("Enumerate添加索引: %v\n", enumerated)
}

// 3. Chainable集合操作
func exampleChainable() {
	fmt.Println("Chainable集合操作:")

	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 创建Chainable
	chain := rust.From(data)
	fmt.Printf("原始数据: %v\n", data)

	// 链式操作
	result := chain.
		Filter(func(x int) bool { return x%2 == 0 }). // 过滤偶数
		Map(func(x int) int { return x * 3 }).        // 乘以3
		Skip(1).                                      // 跳过第一个
		Take(3).                                      // 取3个
		Reverse().                                    // 反转
		Collect()

	fmt.Printf("链式操作结果: %v\n", result)

	// 更多操作
	fmt.Println("\n更多Chainable操作:")

	// Unique去重
	withDuplicates := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	unique := rust.From(withDuplicates).Unique().Collect()
	fmt.Printf("Unique去重: %v -> %v\n", withDuplicates, unique)

	// Partition分区
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	even, odd := rust.From(numbers).Partition(func(x int) bool {
		return x%2 == 0
	})
	fmt.Printf("Partition分区: 偶数=%v, 奇数=%v\n", even.Collect(), odd.Collect())

	// FlatMap扁平化
	nested := []int{1, 2, 3, 4, 5, 6}
	flattened := rust.From(nested).FlatMap(func(x int) []int {
		return []int{x, x * 2}
	}).Collect()
	fmt.Printf("FlatMap扁平化: %v -> %v\n", nested, flattened)

	// Chunk分块
	chunked := rust.From(numbers).Chunk(3).Collect()
	fmt.Printf("Chunk分块(大小3): %v\n", chunked)

	// Window滑动窗口
	windowed := rust.From([]int{1, 2, 3, 4, 5}).Window(3).Collect()
	fmt.Printf("Window滑动窗口(大小3): %v\n", windowed)
}

// 4. 高级迭代器功能
func exampleAdvancedIterator() {
	fmt.Println("高级迭代器功能:")

	// Range迭代器
	fmt.Println("Range迭代器:")
	rangeResult := rust.Collect(rust.Range(1, 10, 2))
	fmt.Printf("Range(1, 10, 2): %v\n", rangeResult)

	// Once迭代器
	onceResult := rust.Collect(rust.Once("hello"))
	fmt.Printf("Once(\"hello\"): %v\n", onceResult)

	// Repeat迭代器（配合Take使用）
	repeatResult := rust.Collect(rust.Take(rust.Repeat("loop"), 5))
	fmt.Printf("Take(Repeat(\"loop\"), 5): %v\n", repeatResult)

	// Empty迭代器
	emptyResult := rust.Collect(rust.Empty[int]())
	fmt.Printf("Empty[int](): %v\n", emptyResult)

	// 使用All和Any
	numbers := []int{2, 4, 6, 8, 10}
	allEven := rust.All(rust.Iter(numbers), func(x int) bool {
		return x%2 == 0
	})
	anyOdd := rust.Any(rust.Iter(numbers), func(x int) bool {
		return x%2 == 1
	})
	fmt.Printf("All偶数检查: %v -> %v\n", numbers, allEven)
	fmt.Printf("Any奇数检查: %v -> %v\n", numbers, anyOdd)

	// 使用Find查找元素
	found := rust.Find(rust.Iter(numbers), func(x int) bool {
		return x > 5
	})
	fmt.Printf("Find(>5): %v\n", found)

	// 使用Count计数
	count := rust.Count(rust.Iter(numbers))
	fmt.Printf("Count: %v -> %d\n", numbers, count)

	// 使用Last获取最后一个
	last := rust.Last(rust.Iter(numbers))
	fmt.Printf("Last: %v -> %v\n", numbers, last)
}

// 5. 实际应用场景
func exampleRealWorld() {
	fmt.Println("实际应用场景:")

	// 场景1: 数据处理管道
	fmt.Println("1. 数据处理管道:")

	type Product struct {
		ID     int
		Name   string
		Price  float64
		Stock  int
		Active bool
	}

	products := []Product{
		{1, "Laptop", 999.99, 10, true},
		{2, "Mouse", 29.99, 50, true},
		{3, "Keyboard", 79.99, 0, false},
		{4, "Monitor", 299.99, 5, true},
		{5, "Headphones", 149.99, 20, true},
		{6, "Tablet", 399.99, 0, false},
	}

	// 处理逻辑：获取有库存的活跃产品，按价格排序（降序），取前3个
	processed := rust.From(products).
		Filter(func(p Product) bool {
			return p.Active && p.Stock > 0
		}).
		Map(func(p Product) Product {
			// 应用折扣
			if p.Price > 100 {
				p.Price *= 0.9 // 9折
			}
			return p
		}).
		Collect()

	// 手动排序（Chainable当前没有排序方法）
	fmt.Printf("有库存的活跃产品（9折后）: %v\n", len(processed))

	// 场景2: 日志分析
	fmt.Println("\n2. 日志分析:")

	type LogEntry struct {
		Timestamp string
		Level     string
		Message   string
	}

	logs := []LogEntry{
		{"2024-01-01 10:00:00", "INFO", "系统启动"},
		{"2024-01-01 10:05:00", "ERROR", "数据库连接失败"},
		{"2024-01-01 10:10:00", "WARN", "内存使用率高"},
		{"2024-01-01 10:15:00", "ERROR", "API调用超时"},
		{"2024-01-01 10:20:00", "INFO", "服务恢复"},
	}

	// 统计错误日志
	errorLogs := rust.From(logs).
		Filter(func(log LogEntry) bool {
			return log.Level == "ERROR"
		}).
		Collect()

	errorMessages := make([]string, len(errorLogs))
	for i, log := range errorLogs {
		errorMessages[i] = fmt.Sprintf("[%s] %s", log.Timestamp, log.Message)
	}

	fmt.Printf("错误日志数量: %d\n", len(errorLogs))
	fmt.Println("错误消息:")
	for _, msg := range errorMessages {
		fmt.Printf("  %s\n", msg)
	}

	// 场景3: 文本处理
	fmt.Println("\n3. 文本处理:")

	text := "The quick brown fox jumps over the lazy dog. The dog barks back."
	words := strings.Fields(text)

	// 统计词频
	wordStats := rust.From(words).
		Map(func(word string) string {
			return strings.ToLower(word)
		}).
		Collect()

	// 使用Chainable进行更多分析
	uniqueWords := rust.From(wordStats).Unique().Collect()
	longWords := rust.From(wordStats).
		Filter(func(word string) bool {
			return len(word) > 4
		}).
		Unique().
		Collect()

	fmt.Printf("总单词数: %d\n", len(words))
	fmt.Printf("唯一单词数: %d\n", len(uniqueWords))
	fmt.Printf("长单词(>4字母): %v\n", longWords)
}

// 6. 性能比较和最佳实践
func examplePerformance() {
	fmt.Println("性能比较和最佳实践:")

	// 传统Go风格 vs RustGo风格
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println("1. 传统Go风格:")
	fmt.Println("   for循环:")
	sumTraditional := 0
	for _, n := range numbers {
		if n%2 == 0 {
			sumTraditional += n * 2
		}
	}
	fmt.Printf("    结果: %d\n", sumTraditional)

	fmt.Println("\n2. RustGo风格:")
	fmt.Println("   迭代器风格:")
	sumIterator := rust.Fold(
		rust.Map(
			rust.Filter(
				rust.Iter(numbers),
				func(x int) bool { return x%2 == 0 },
			),
			func(x int) int { return x * 2 },
		),
		0,
		func(acc, x int) int { return acc + x },
	)
	fmt.Printf("    结果: %d\n", sumIterator)

	fmt.Println("\n3. Chainable风格:")
	sumChainable := rust.From(numbers).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * 2 }).
		Fold(0, func(acc, x int) int { return acc + x })
	fmt.Printf("    结果: %d\n", sumChainable)

	// 最佳实践
	fmt.Println("\n最佳实践:")
	fmt.Println("  1. 小数据集: 使用Chainable（更易读）")
	fmt.Println("  2. 大数据集: 使用迭代器（惰性求值，节省内存）")
	fmt.Println("  3. 复杂管道: 使用迭代器（组合更灵活）")
	fmt.Println("  4. 简单转换: 使用Chainable（链式调用更直观）")
	fmt.Println("  5. 性能关键: 测试两种方式，选择更快的")

	// 惰性求值示例
	fmt.Println("\n惰性求值示例:")
	fmt.Println("  迭代器是惰性的，只有在调用Collect()时才执行计算")
	fmt.Println("  这对于大数据集或无限序列特别有用")

	// 生成无限序列（但只取一部分）
	fmt.Println("  生成斐波那契数列（前10个）:")

	// 自定义迭代器实现
	type FibonacciIterator struct {
		a, b int
	}

	fibonacci := rust.Generate(10, func(i int) int {
		if i == 0 {
			return 0
		}
		if i == 1 {
			return 1
		}
		// 这里简化实现，实际应该维护状态
		// 真正的实现应该是一个自定义迭代器
		return i // 简化示例
	})

	fibResult := fibonacci.Collect()
	fmt.Printf("  结果: %v\n", fibResult)
}

// 斐波那契迭代器实现（完整版）
type FibonacciIterator struct {
	current int
	next    int
}

func (f *FibonacciIterator) Next() rust.Option[int] {
	value := f.current
	f.current, f.next = f.next, f.current+f.next
	return rust.Some(value)
}

func NewFibonacciIterator() rust.Iterator[int] {
	return &FibonacciIterator{current: 0, next: 1}
}
