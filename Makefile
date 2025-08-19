.PHONY: all

all: run test clean
# Цель для запуска приложения
run:
	go run main.go

# Цель для запуска тестов
test:
	go test ./...

# Цель для очистки временных файлов
clean:
	rm -f ./myapp

# Цель для обновления модулей
tidy:
	go mod tidy