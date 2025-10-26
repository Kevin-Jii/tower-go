# Makefile
.PHONY: run docs

# 每次运行前确保文档是最新的
run: docs
	go run cmd/main.go

# 专门的文档生成命令
docs:
	# 运行 swag init 生成 docs 文件夹
	swag init -g cmd/main.go

build: docs
	go build -o server cmd/main.go