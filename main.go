package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Apply the CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/uploads", "./uploads")

	// Set a POST route to handle file uploads
	router.POST("/upload", func(c *gin.Context) {
		// Single file
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		// Check file extension
		ext := filepath.Ext(file.Filename)
		if ext != ".mp3" {
			c.String(http.StatusBadRequest, "only .mp3 files are allowed")
			return
		}

		// Save the file to a specified directory
		if err := c.SaveUploadedFile(file, filepath.Join("uploads", file.Filename)); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// Set a GET route to list uploaded files
	router.GET("/files", func(c *gin.Context) {
		files, err := listFiles("uploads")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("could not list files: %s", err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": files})
	})

	// Create the uploads directory if not exists
	if err := createUploadsDir(); err != nil {
		log.Fatalf("Could not create uploads directory: %s", err)
	}

	// Run the server
	router.Run(":8080")
}

// createUploadsDir creates the uploads directory if it does not exist
func createUploadsDir() error {
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// listFiles lists all files in the specified directory
func listFiles(directory string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mp3") {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}
