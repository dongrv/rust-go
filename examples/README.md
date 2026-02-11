# RustGo 示例目录

本目录包含 RustGo 库的完整使用示例，展示了如何在 Go 项目中应用 Rust 风格的编程模式。

## 示例文件说明

### 1. `option_example.go` - Option 类型示例
展示了如何使用 Rust 风格的 `Option[T]` 类型进行空安全编程。

**主要内容：**
- 基本创建和使用：`Some()` 和 `None()`
- 安全解包：`UnwrapOr()`, `UnwrapOrElse()`, `Expect()`
- 链式操作：`MapOption()`, `AndThenOption()`
- 实际应用场景：数据库查询、配置解析
- 组合使用和与传统 Go 风格的比较

**运行方式：**
```bash
cd examples
go run option_example.go
```

### 2. `result_example.go` - Result 类型示例
展示了如何使用 Rust 风格的 `Result[T, E]` 类型进行错误处理。

**主要内容：**
- 基本创建和使用：`Ok()` 和 `Err()`
- 错误处理模式：`UnwrapOr()`, `UnwrapOrElse()`, `Or()`, `OrElse()`
- 链式操作（铁路编程）：`AndThenResult()`, `MapResult()`
- 实际应用场景：文件操作、API 调用
- 高级用法：错误恢复、组合模式

**运行方式：**
```bash
cd examples
go run result_example.go
```

### 3. `iterator_example.go` - Iterator 和 Chainable 示例
展示了如何使用 Rust 风格的迭代器和链式集合操作。

**主要内容：**
- 基本迭代器操作：`Iter()`, `Collect()`, `ForEach()`, `Fold()`, `Reduce()`
- 链式迭代器操作：`Map()`, `Filter()`, `Take()`, `Skip()`, `Chain()`, `Zip()`, `Enumerate()`
- Chainable 集合操作：`From()`, 链式方法调用
- 高级迭代器功能：`Range()`, `Once()`, `Repeat()`, `All()`, `Any()`, `Find()`
- 实际应用场景：数据处理管道、日志分析、文本处理
- 性能比较和最佳实践

**运行方式：**
```bash
cd examples
go run iterator_example.go
```

### 4. `comprehensive_example.go` - 综合示例
展示了一个完整的业务场景：用户订单处理系统，综合使用所有 RustGo 特性。

**主要内容：**
- 完整的业务类型定义：`User`, `Product`, `Order`, `OrderItem`
- 使用 `Option` 处理可能为空的值
- 使用 `Result` 处理可能失败的操作
- 使用链式操作处理业务逻辑
- 完整的订单创建流程（铁路编程模式）
- 使用迭代器处理批量操作
- 查询和分析功能
- 错误处理和恢复机制
- 复杂的数据转换和分析

**业务场景：**
模拟一个电商平台的订单处理系统，包括：
1. 用户和产品管理
2. 订单验证和创建
3. 库存管理
4. 数据分析和统计
5. 错误处理和重试机制

**运行方式：**
```bash
cd examples
go run comprehensive_example.go
```

## 如何运行示例

### 前提条件
1. 确保已安装 Go 1.21 或更高版本
2. 确保在项目根目录下

### 运行单个示例
```bash
cd examples
go run option_example.go
```

### 运行所有示例
```bash
cd examples
for file in *.go; do
    echo "=== 运行 $file ==="
    go run "$file"
    echo
done
```

## 示例设计原则

### 1. 渐进式学习
- 从简单到复杂
- 每个示例专注于一个核心概念
- 综合示例展示实际应用

### 2. 实用性
- 基于真实业务场景
- 展示最佳实践
- 提供性能建议

### 3. 可测试性
- 每个示例都是独立的
- 包含完整的错误处理
- 输出清晰的结果

### 4. 教育性
- 包含详细的注释
- 展示与传统 Go 风格的对比
- 解释设计决策

## 核心概念展示

### Option 模式
```go
// 传统 Go 风格
var value *int
if value != nil {
    // 使用 value
}

// RustGo 风格
value := rustgo.Some(42)
value.ForEach(func(x int) {
    // 安全地使用 value
})
```

### Result 模式（铁路编程）
```go
// 传统 Go 风格
result1, err := step1()
if err != nil {
    return err
}
result2, err := step2(result1)
if err != nil {
    return err
}
// ...

// RustGo 风格
finalResult := rustgo.AndThenResult(
    step1(),
    func(r1) rustgo.Result { return step2(r1) },
    // ...
)
```

### 迭代器和 Chainable
```go
// 传统 Go 风格
var result []int
for _, x := range data {
    if x%2 == 0 {
        result = append(result, x*2)
    }
}

// RustGo 风格
result := rustgo.From(data).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * 2 }).
    Collect()
```

## 最佳实践建议

### 1. 何时使用 Option
- 替代可能为 `nil` 的指针
- 函数可能返回空值的情况
- 配置项或可选参数

### 2. 何时使用 Result
- 函数可能失败的情况
- 需要明确错误类型的场景
- 复杂的错误处理链

### 3. 何时使用 Iterator
- 处理大数据集（惰性求值）
- 需要组合多个操作
- 处理无限序列

### 4. 何时使用 Chainable
- 小到中等数据集
- 需要链式调用的场景
- 代码可读性更重要时

## 扩展示例

您可以根据这些示例创建自己的应用：
1. **Web 服务**：使用 `Result` 处理 HTTP 请求错误
2. **数据处理**：使用 `Iterator` 处理流式数据
3. **配置管理**：使用 `Option` 处理可选配置
4. **数据库操作**：组合使用所有类型构建健壮的 DAO 层

## 贡献

欢迎贡献更多示例！请确保：
1. 示例是完整且可运行的
2. 包含详细的注释
3. 遵循项目的代码风格
4. 展示最佳实践

## 许可证

这些示例遵循与 RustGo 库相同的 MIT 许可证。

## 更多资源

- [RustGo 主文档](../README.md)
- [Go 官方文档](https://golang.org/doc/)
- [Rust 官方文档](https://doc.rust-lang.org/book/)

---

通过这些示例，您将能够：
1. 理解 RustGo 的核心概念
2. 在实际项目中使用 Rust 风格的编程模式
3. 编写更安全、更表达力强的 Go 代码
4. 构建健壮且易于维护的应用程序

Happy coding! 🦀 + 🐹 = ❤️