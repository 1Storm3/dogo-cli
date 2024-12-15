package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dogo",
		Short: "dogo: Cli утилита для генерации сервисов с REST и gRPC",
		Long:  `dogo: Cli утилита для генерации сервисов с REST и gRPC`,
	}
	rootCmd.SetHelpCommand(nil)

	rootCmd.AddCommand(createGenerateCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createGenerateCommand() *cobra.Command {
	var serviceName string
	var serviceType string

	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Генерация микросервиса",
		RunE: func(cmd *cobra.Command, args []string) error {
			if serviceName == "" {
				return fmt.Errorf("Сервисное имя не может быть пустым")
			}

			if serviceType != "rest" && serviceType != "grpc" {
				return fmt.Errorf("Тип сервиса может быть только 'rest' или 'grpc'")
			}

			return generateService(serviceName, serviceType)
		},
	}

	cmd.Flags().StringVarP(&serviceName, "name", "n", "", "Имя микросервиса (required)")
	cmd.Flags().StringVarP(&serviceType, "type", "t", "rest", "Тип сервиса: rest or grpc (default: rest)")

	return cmd
}

func generateService(serviceName, serviceType string) error {
	fmt.Printf("Генерация %s микросервиса: %s\n", serviceType, serviceName)

	dirStructure := map[string][]string{
		"dogo":     {serviceName},
		"config":   {"config"},
		"internal": {"router", "handler", "service", "repository", "middleware"},
		"pkg":      {"utils"},
		"api":      {"proto"},
		"docs":     {"swagger"},
	}

	if err := createDirectories(serviceName, dirStructure); err != nil {
		return err
	}

	if err := createMainFile(serviceName); err != nil {
		return err
	}

	if err := createAdditionalFiles(serviceName); err != nil {
		return err
	}

	if err := initializeGoMod(serviceName); err != nil {
		return err
	}

	fmt.Println("Микросервис успешно генерирован!")
	return nil
}

func initializeGoMod(serviceName string) error {
	fmt.Println("Инициализация go.mod...")
	servicePath := filepath.Join(serviceName)

	cmd := exec.Command("go", "mod", "init", serviceName)
	cmd.Dir = servicePath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Ошибка при инициализации go.mod: %v, вывод: %s", err, output)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = servicePath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Ошибка при установке зависимостей: %v, вывод: %s", err, output)
	}

	return nil
}

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
