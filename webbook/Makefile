# 定义一个伪目标，确保在调用make docker的时候不会执行相关文件
.PHONY: docker
# 名为docker的目标

build:
	if exist webook.exe del webook.exe
	set GOOS=linux
	set GOARCH=amd64
	go build -o webook .

docker:
#@ 标识执行命令的时候会打印命令的内容
	@del webook || true
	@GOOS=linux GOARCH=amd64 go build -o webook .
	@docker rmi -f webook:latest || true
	@docker build -t webook:v0.0.1 .

