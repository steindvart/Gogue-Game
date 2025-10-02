# Сборка всех бинарников
.PHONY: build
build: prepare-bin
# 	go build -o bin/t01_goroutine_clocks.exe ./t01_goroutine_clocks

TEST_FLAGS=-v

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

# Запуск всех тестов со сбором покрытия
COV_PATH=bin/coverage.out
COV_HTML_PATH=bin/coverage.html
TEST_FLAGS_COV=$(TEST_FLAGS) -cover -coverprofile $(COV_PATH)

.PHONY: test-coverage
test-coverage: prepare-bin
	go test $(TEST_FLAGS_COV) ./... && go tool cover -func $(COV_PATH)
	go tool cover -html $(COV_PATH) -o $(COV_HTML_PATH)

.PHONY: prepare-bin
prepare-bin:
	@if [ ! -d "bin" ]; then \
			echo "Makefile - Make directory 'bin/'"; \
			mkdir -p bin; \
	fi

# Запуск линтера (если установлен golangci-lint)
.PHONY: lint
lint:
	golangci-lint run ./...

# Очистка артефактов сборки
.PHONY: clean
clean:
	rm -rf bin/*
	go clean -testcache

# Сборка и запуск отдельного приложения
# .PHONY: run-t01
# run-t01:
# 	go run ./t01_goroutine_clocks/main.go 10 1000
