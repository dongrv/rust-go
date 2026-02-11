// 综合示例 - 展示RustGo库的完整用法
// 这个示例展示了如何在实际项目中使用Option、Result、Iterator和Chainable

package main

import (
	"fmt"
	"time"

	"github.com/dongrv/rust-go"
)

// 模拟一个完整的业务场景：用户订单处理系统

// 定义业务类型
type UserID int
type ProductID int
type OrderID string

type User struct {
	ID       UserID
	Name     string
	Email    string
	IsActive bool
}

type Product struct {
	ID       ProductID
	Name     string
	Price    float64
	Stock    int
	Category string
}

type OrderItem struct {
	ProductID ProductID
	Quantity  int
	Price     float64
}

type Order struct {
	ID        OrderID
	UserID    UserID
	Items     []OrderItem
	Total     float64
	Status    string
	CreatedAt time.Time
}

// 模拟数据库/存储层
type Database struct {
	users    map[UserID]User
	products map[ProductID]Product
	orders   map[OrderID]Order
}

// 业务逻辑层使用RustGo类型
type OrderService struct {
	db *Database
}

// 创建新的订单服务
func NewOrderService(db *Database) *OrderService {
	return &OrderService{db: db}
}

// 1. 使用Option处理可能为空的值
func (s *OrderService) FindUserByID(id UserID) rust.Option[User] {
	user, exists := s.db.users[id]
	if !exists {
		return rust.None[User]()
	}
	return rust.Some(user)
}

func (s *OrderService) FindProductByID(id ProductID) rust.Option[Product] {
	product, exists := s.db.products[id]
	if !exists {
		return rust.None[Product]()
	}
	return rust.Some(product)
}

// 2. 使用Result处理可能失败的操作
func (s *OrderService) ValidateOrderRequest(
	userID UserID,
	items []OrderItem,
) rust.Result[[]OrderItem, string] {
	// 验证用户存在且活跃
	userOpt := s.FindUserByID(userID)
	if userOpt.IsNone() {
		return rust.Err[[]OrderItem, string]("用户不存在")
	}

	user := userOpt.Unwrap()
	if !user.IsActive {
		return rust.Err[[]OrderItem, string]("用户账户未激活")
	}

	// 验证每个订单项
	validatedItems := []OrderItem{}
	for _, item := range items {
		productOpt := s.FindProductByID(item.ProductID)
		if productOpt.IsNone() {
			return rust.Err[[]OrderItem, string](
				fmt.Sprintf("产品ID %d 不存在", item.ProductID),
			)
		}

		product := productOpt.Unwrap()
		if product.Stock < item.Quantity {
			return rust.Err[[]OrderItem, string](
				fmt.Sprintf("产品 %s 库存不足", product.Name),
			)
		}

		if item.Quantity <= 0 {
			return rust.Err[[]OrderItem, string]("购买数量必须大于0")
		}

		// 使用产品当前价格
		validatedItem := OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		validatedItems = append(validatedItems, validatedItem)
	}

	return rust.Ok[[]OrderItem, string](validatedItems)
}

// 3. 使用链式操作处理业务逻辑
func (s *OrderService) CalculateOrderTotal(items []OrderItem) float64 {
	// 使用迭代器计算订单总额
	total := rust.Fold(rust.Iter(items), 0.0, func(acc float64, item OrderItem) float64 {
		return acc + (item.Price * float64(item.Quantity))
	})

	return total
}

// 4. 完整的订单创建流程（铁路编程模式）
func (s *OrderService) CreateOrder(
	userID UserID,
	items []OrderItem,
) rust.Result[Order, string] {
	// 使用铁路编程模式：一系列操作，任何一个失败都会短路
	return rust.AndThenResult(
		s.ValidateOrderRequest(userID, items),
		func(validatedItems []OrderItem) rust.Result[Order, string] {
			// 计算订单总额
			total := s.CalculateOrderTotal(validatedItems)

			// 生成订单ID
			orderID := OrderID(fmt.Sprintf("ORD-%d-%d", userID, time.Now().Unix()))

			// 创建订单
			order := Order{
				ID:        orderID,
				UserID:    userID,
				Items:     validatedItems,
				Total:     total,
				Status:    "pending",
				CreatedAt: time.Now(),
			}

			// 保存订单（模拟）
			s.db.orders[orderID] = order

			// 更新库存
			updateResult := s.UpdateInventory(validatedItems)
			if updateResult.IsErr() {
				// 回滚订单创建
				delete(s.db.orders, orderID)
				return rust.Err[Order, string](updateResult.UnwrapErr())
			}

			return rust.Ok[Order, string](order)
		},
	)
}

// 5. 使用迭代器处理批量操作
func (s *OrderService) UpdateInventory(items []OrderItem) rust.Result[bool, string] {
	// 使用迭代器处理每个产品的库存更新
	iter := rust.Iter(items)
	for {
		next := iter.Next()
		if next.IsNone() {
			break
		}

		item := next.Unwrap()
		productOpt := s.FindProductByID(item.ProductID)
		if productOpt.IsNone() {
			return rust.Err[bool, string](
				fmt.Sprintf("产品ID %d 不存在", item.ProductID),
			)
		}

		product := productOpt.Unwrap()
		if product.Stock < item.Quantity {
			return rust.Err[bool, string](
				fmt.Sprintf("产品 %s 库存不足，当前库存: %d", product.Name, product.Stock),
			)
		}

		// 更新库存
		product.Stock -= item.Quantity
		s.db.products[item.ProductID] = product
	}

	return rust.Ok[bool, string](true)
}

// 6. 查询和分析功能
func (s *OrderService) GetUserOrders(userID UserID) []Order {
	// 使用Chainable过滤和转换数据
	userOrders := FromMapValues(s.db.orders).
		Filter(func(order Order) bool {
			return order.UserID == userID
		}).
		Collect()

	return userOrders
}

func (s *OrderService) GetTopProducts(limit int) []Product {
	// 统计产品销量
	productSales := make(map[ProductID]int)

	// 使用迭代器遍历所有订单
	rust.ForEach(IterMapValues(s.db.orders), func(order Order) {
		rust.ForEach(rust.Iter(order.Items), func(item OrderItem) {
			productSales[item.ProductID] += item.Quantity
		})
	})

	// 转换为切片并排序
	var products []Product
	for pair := range productSales {
		product, exists := s.db.products[pair]
		if exists && product.Stock > 0 {
			products = append(products, product)
		}
	}

	// 使用Chainable进行过滤和限制
	result := rust.From(products).
		Filter(func(product Product) bool {
			return product.Stock > 0
		}).
		Take(limit).
		Collect()

	return result
}

// 7. 错误处理和恢复
func (s *OrderService) ProcessOrderWithRetry(
	userID UserID,
	items []OrderItem,
	maxRetries int,
) rust.Result[Order, string] {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		result := s.CreateOrder(userID, items)
		if result.IsOk() {
			return result
		}

		if attempt < maxRetries {
			fmt.Printf("订单创建失败（尝试 %d/%d）: %v\n",
				attempt, maxRetries, result.UnwrapErr())
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	return rust.Err[Order, string]("订单创建失败，已达到最大重试次数")
}

// 辅助函数：从map创建Chainable
func FromMapValues[K comparable, V any](m map[K]V) *rust.Chainable[V] {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return rust.From(values)
}

func IterMapValues[K comparable, V any](m map[K]V) rust.Iterator[V] {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return rust.Iter(values)
}

func FromMap[K comparable, V any](m map[K]V) *rust.Chainable[rust.Pair[K, V]] {
	pairs := make([]rust.Pair[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, rust.Pair[K, V]{First: k, Second: v})
	}
	return rust.From(pairs)
}

// 主函数：演示完整的使用场景
func main() {
	fmt.Println("=== RustGo综合示例：订单处理系统 ===")
	fmt.Println()

	// 初始化数据库
	db := &Database{
		users: map[UserID]User{
			1: {ID: 1, Name: "张三", Email: "zhangsan@example.com", IsActive: true},
			2: {ID: 2, Name: "李四", Email: "lisi@example.com", IsActive: false},
			3: {ID: 3, Name: "王五", Email: "wangwu@example.com", IsActive: true},
		},
		products: map[ProductID]Product{
			101: {ID: 101, Name: "笔记本电脑", Price: 5999.99, Stock: 10, Category: "电子产品"},
			102: {ID: 102, Name: "无线鼠标", Price: 199.99, Stock: 50, Category: "电子产品"},
			103: {ID: 103, Name: "机械键盘", Price: 499.99, Stock: 20, Category: "电子产品"},
			104: {ID: 104, Name: "显示器", Price: 1299.99, Stock: 5, Category: "电子产品"},
			105: {ID: 105, Name: "办公椅", Price: 899.99, Stock: 0, Category: "家具"},
		},
		orders: make(map[OrderID]Order),
	}

	// 创建订单服务
	service := NewOrderService(db)

	// 演示1: 使用Option处理可能为空的值
	fmt.Println("1. 使用Option处理可能为空的值:")
	userOpt := service.FindUserByID(1)
	if userOpt.IsSome() {
		user := userOpt.Unwrap()
		fmt.Printf("  找到用户: %s (%s)\n", user.Name, user.Email)
	}

	noneUserOpt := service.FindUserByID(999)
	if noneUserOpt.IsNone() {
		fmt.Println("  用户不存在时返回None")
	}

	// 使用UnwrapOr提供默认值
	userName := rust.MapOption(userOpt, func(u User) string { return u.Name }).
		UnwrapOr("未知用户")
	fmt.Printf("  用户名（使用UnwrapOr）: %s\n", userName)

	// 演示2: 使用Result处理可能失败的操作
	fmt.Println("\n2. 使用Result处理可能失败的操作:")

	// 有效的订单请求
	validItems := []OrderItem{
		{ProductID: 101, Quantity: 1},
		{ProductID: 102, Quantity: 2},
	}

	validationResult := service.ValidateOrderRequest(1, validItems)
	if validationResult.IsOk() {
		fmt.Println("  订单验证成功")
	} else {
		fmt.Printf("  订单验证失败: %v\n", validationResult.UnwrapErr())
	}

	// 无效的订单请求（用户未激活）
	invalidUserItems := []OrderItem{{ProductID: 101, Quantity: 1}}
	invalidResult := service.ValidateOrderRequest(2, invalidUserItems)
	if invalidResult.IsErr() {
		fmt.Printf("  用户未激活错误: %v\n", invalidResult.UnwrapErr())
	}

	// 演示3: 完整的订单创建流程
	fmt.Println("\n3. 完整的订单创建流程（铁路编程）:")

	orderResult := service.CreateOrder(1, validItems)
	if orderResult.IsOk() {
		order := orderResult.Unwrap()
		fmt.Printf("  订单创建成功!\n")
		fmt.Printf("  订单ID: %s\n", order.ID)
		fmt.Printf("  订单总额: ¥%.2f\n", order.Total)
		fmt.Printf("  订单状态: %s\n", order.Status)
		fmt.Printf("  创建时间: %s\n", order.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	if orderResult.IsErr() {
		fmt.Printf("  订单创建失败: %v\n", orderResult.UnwrapErr())
	}

	// 演示4: 使用Chainable进行数据分析
	fmt.Println("\n4. 使用Chainable进行数据分析:")

	// 获取用户订单
	userOrders := service.GetUserOrders(1)
	fmt.Printf("  用户1的订单数量: %d\n", len(userOrders))

	// 获取热销产品
	topProducts := service.GetTopProducts(3)
	fmt.Println("  热销产品（有库存）:")
	rust.From(topProducts).ForEach(func(product Product) {
		fmt.Printf("    - %s (¥%.2f, 库存: %d)\n",
			product.Name, product.Price, product.Stock)
	})

	// 演示5: 错误处理和重试机制
	fmt.Println("\n5. 错误处理和重试机制:")

	// 尝试购买缺货的产品
	outOfStockItems := []OrderItem{
		{ProductID: 105, Quantity: 1}, // 办公椅库存为0
	}

	retryResult := service.ProcessOrderWithRetry(1, outOfStockItems, 3)
	if retryResult.IsErr() {
		fmt.Printf("  重试后仍然失败: %v\n", retryResult.UnwrapErr())
	}

	// 演示6: 复杂的业务逻辑组合
	fmt.Println("\n6. 复杂的业务逻辑组合:")

	// 批量处理多个订单请求
	orderRequests := []struct {
		userID UserID
		items  []OrderItem
	}{
		{1, []OrderItem{{ProductID: 103, Quantity: 1}}},
		{3, []OrderItem{{ProductID: 102, Quantity: 3}}},
		{999, []OrderItem{{ProductID: 101, Quantity: 1}}}, // 无效用户
		{1, []OrderItem{{ProductID: 105, Quantity: 1}}},   // 缺货产品
	}

	fmt.Println("  批量处理订单请求:")
	for i, req := range orderRequests {
		result := service.CreateOrder(req.userID, req.items)
		status := "成功"
		if result.IsErr() {
			status = result.UnwrapErr()
		}
		fmt.Printf("  请求%d: 用户%d -> %s\n", i+1, req.userID, status)
	}

	// 演示7: 使用迭代器进行复杂的数据转换
	fmt.Println("\n7. 使用迭代器进行复杂的数据转换:")

	// 计算所有订单的总销售额
	allOrders := FromMapValues(db.orders)
	totalRevenue := rust.Fold(allOrders.Iter(), 0.0, func(acc float64, order Order) float64 {
		return acc + order.Total
	})

	fmt.Printf("  所有订单总销售额: ¥%.2f\n", totalRevenue)

	// 按用户统计订单数量
	userOrderCounts := make(map[UserID]int)
	rust.ForEach(allOrders.Iter(), func(order Order) {
		userOrderCounts[order.UserID]++
	})

	fmt.Println("  用户订单统计:")
	for userID, count := range userOrderCounts {
		userOpt := service.FindUserByID(userID)
		userName := rust.MapOption(userOpt, func(u User) string { return u.Name }).
			UnwrapOr("未知用户")
		fmt.Printf("    %s: %d 个订单\n", userName, count)
	}

	// 总结
	fmt.Println("\n=== 总结 ===")
	fmt.Println("通过这个综合示例，我们展示了:")
	fmt.Println("1. 使用Option替代nil指针，避免空指针异常")
	fmt.Println("2. 使用Result进行明确的错误处理，支持铁路编程")
	fmt.Println("3. 使用Iterator进行惰性求值，适合大数据集")
	fmt.Println("4. 使用Chainable进行链式操作，代码更易读")
	fmt.Println("5. 组合使用这些类型构建健壮的业务逻辑")
	fmt.Println("6. 错误恢复和重试机制")
	fmt.Println("7. 复杂的数据分析和转换")
	fmt.Println()
	fmt.Println("RustGo让Go代码更安全、更表达力强、更易于维护！")
}

// Chainable已经包含ForEach方法，这里不需要重复定义
// Result类型没有ForEach方法，我们使用IsOk()和Unwrap()来访问值
