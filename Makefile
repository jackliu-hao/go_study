.PHONY: mockSvc


mockSvc:
	@mockgen -source=webbook/internal/service/user.go -package=svcmocks -destination=webbook/internal/service/mocks/user_mock_gen.go
	@mockgen -source=webbook/internal/service/code.go -package=svcmocks -destination=webbook/internal/service/mocks/code_mock_gen.go
	@go mod tidy

mockRepo:
	@mockgen -source=webbook/internal/repository/user.go -package=repomocks -destination=webbook/internal/repository/mocks/user_mock_gen.go
	@mockgen -source=webbook/internal/repository/code.go -package=repomocks -destination=webbook/internal/repository/mocks/code_mock_gen.go
	@go mod tidy

mockDao:
	@mockgen -source=webbook/internal/repository/dao/user.go -package=daomocks -destination=webbook/internal/repository/dao/mocks/user_mock_gen.go
	@mockgen -source=webbook/internal/repository/cache/user.go -package=daomocks -destination=webbook/internal/repository/cache/mocks/user_mock_gen.go
	@go mod tidy

mockRedis:
# 生成redis的mock , 因为是使用三方的包，这里不存在--source， 第三个参数是对应的包路径和类名，可以使用 , 生成多个类型
	@mockgen -package=redismocks -destination=webbook/mocks/redis_mock_gen.go  github.com/redis/go-redis/v9 Cmdable
	@go mod tidy

mockRateLimit:
	@mockgen -source=webbook/pkg/ratelimit/type.go -package=limitmocks -destination=webook/pkg/ratelimit/mocks/limiter.mock.go

mockAll:
	#调用所有的mock命令
	@$(MAKE) mockSvc
	@$(MAKE) mockRepo
	@$(MAKE) mockDao
	@$(MAKE) mockRedis
	@$(MAKE) mockRateLimit
	@go mod tidy
	@echo "mock all done"

clean:
	@rm -rf webbook/internal/service/mocks/user_mock_gen.go
	@rm -rf webbook/internal/service/mocks/code_mock_gen.go
	@rm -rf webbook/internal/repository/mocks/user_mock_gen.go
	@rm -rf webbook/internal/repository/mocks/code_mock_gen.go
	@rm -rf webbook/internal/repository/dao/mocks/user_mock_gen.go