.PHONY: mockSvc


mockSvc:
	@mockgen -source=webbook/internal/service/user.go -package=svcmocks -destination=webbook/internal/service/mocks/user_mock_gen.go
	@mockgen -source=webbook/internal/service/code.go -package=svcmocks -destination=webbook/internal/service/mocks/code_mock_gen.go
	@go mod tidy