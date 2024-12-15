.PHONY: build install run deps clean

# Путь к исполняемому файлу
BINARY=./dogo

# Команда для установки зависимостей
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Компиляция утилиты
build: deps
	@echo "Установка зависимостей..."
	go build -o $(BINARY) .

# Установка утилиты глобально
install: build
	@echo "Установка утилиты..."
	cp $(BINARY) /usr/local/bin/dogo
	@echo "Установка завершена. Используй 'dogo --help' для получения справки."

# Запуск утилиты напрямую
run:
	@echo "Запуск..."
	go run main.go

# Очистка сборочных файлов
clean:
	@echo "Очистка сборочных файл..."
	rm -f $(BINARY)
