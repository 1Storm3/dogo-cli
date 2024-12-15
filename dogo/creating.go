package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func createDirectories(base string, structure map[string][]string) error {
	for root, subdirs := range structure {
		rootPath := filepath.Join(base, root)
		if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
			return fmt.Errorf("Ошибка при создании директории %s: %w", rootPath, err)
		}

		for _, sub := range subdirs {
			subPath := filepath.Join(rootPath, sub)
			if err := os.MkdirAll(subPath, os.ModePerm); err != nil {
				return fmt.Errorf("Ошибка при создании директории %s: %w", subPath, err)
			}
		}
	}
	return nil
}

func createMainFile(serviceName string) error {
	mainFilePath := filepath.Join(serviceName, "dogo", serviceName, "main.go")

	mainContent := fmt.Sprintf(`package main

import (
	"log"
	"%s/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting %s service...")

	r := gin.Default()

	handler.SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
`, serviceName, serviceName)

	if err := os.MkdirAll(filepath.Dir(mainFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("Ошибка при создании директории: %w", err)
	}

	if err := os.WriteFile(mainFilePath, []byte(mainContent), os.ModePerm); err != nil {
		return fmt.Errorf("Ошибка при создании файла: %w", err)
	}

	return nil
}

func createAdditionalFiles(serviceName string) error {
	handlerFilePath := filepath.Join(serviceName, "internal", "router", "routes.go")

	handlerContent := `package router

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}`

	if err := os.WriteFile(handlerFilePath, []byte(handlerContent), os.ModePerm); err != nil {
		return fmt.Errorf("Ошибка при создании файла: %w", err)
	}

	return nil
}
