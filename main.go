package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	createFolderResult()

	go func() {
		for {
			cleanUpOldFiles("public/result", 25*time.Minute)
			time.Sleep(2 * time.Minute)
		}
	}()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Static("/", "./public")

	app.Listen(":" + os.Getenv("PORT"))
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error load .env file")
	}
}

func createFolderResult() {
	dir := filepath.Join("public", "result")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Failed to create directory %s: %v", dir, err)
	}
}

func cleanUpOldFiles(dir string, maxAge time.Duration) {
	files, err := os.ReadDir(dir)

	if err != nil {
		log.Printf("Failed to read directory %s: %v", dir, err)
		return
	}

	now := time.Now()
	for _, file := range files {
		info, err := file.Info()

		if err != nil {
			log.Printf("Failed to get file info %s: %v", file.Name(), err)
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			filePath := filepath.Join(dir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Printf("Failed to delete file %s: %v", filePath, err)
			} else {
				log.Printf("File %s deleted successfully", filePath)
			}
		}
	}
}
