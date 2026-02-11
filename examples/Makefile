# RustGo Examples Makefile
# 用于管理和运行示例的Makefile

.PHONY: all option result iterator comprehensive test clean help

# 默认目标：运行所有示例
all: option result iterator comprehensive

# 运行Option示例
option:
	@echo "=== 运行Option示例 ==="
	@go run option_example.go

# 运行Result示例
result:
	@echo "=== 运行Result示例 ==="
	@go run result_example.go

# 运行Iterator示例
iterator:
	@echo "=== 运行Iterator示例 ==="
	@go run iterator_example.go

# 运行综合示例
comprehensive:
	@echo "=== 运行综合示例 ==="
	@go run comprehensive_example.go

# 测试所有示例
test:
	@echo "=== 测试所有示例 ==="
	@./run_all.sh || ./run_all.bat

# 清理生成的文件
clean:
	@echo "清理生成的文件..."
	@rm -f *.exe
	@rm -f *~
	@echo "清理完成"

# 显示帮助信息
help:
	@echo "RustGo 示例管理工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make all          运行所有示例"
	@echo "  make option       运行Option示例"
	@echo "  make result       运行Result示例"
	@echo "  make iterator     运行Iterator示例"
	@echo "  make comprehensive 运行综合示例"
	@echo "
