package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Image represents the model for storing image information in the database.
type Image struct {
	gorm.Model
	Name string
}

var db *gorm.DB

func main() {
	// Initialize the database connection.
	initDB()
	defer db.Close()

	// Create a Gin router.
	router := gin.Default()
	router.Use(corsMiddleware())

	// Serve static files (for serving images).
	router.Static("/uploads", "./uploads")

	// Define routes.
	router.POST("/upload", uploadImage)
	router.GET("/images/:id", getImage)

	// Start the server.
	router.Run(":8080")
}

// Initialize the database connection using GORM.
func initDB() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=admin dbname=images password=psltest sslmode=disable")
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto-migrate the schema to create the 'images' table.
	db.AutoMigrate(&Image{})
}

// Middleware to allow Cross-Origin Resource Sharing (CORS).
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// Handle image upload.
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read the uploaded file into memory.
	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileBytes.Close()

	// Create a new Image record and save it to the database.
	image := Image{Name: file.Filename}
	if err := db.Create(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the uploaded file data to the 'uploads' folder with the image ID as the filename.
	// Note: You should handle file naming and storage securely in a real-world application.
	// Here, we use the ID as the filename for simplicity.
	// Make sure the 'uploads' folder exists.
	err = c.SaveUploadedFile(file, "uploads/"+fmt.Sprintf("%d", image.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "Id": image.ID})
}

// Handle image retrieval.
func getImage(c *gin.Context) {
	id := c.Param("id")

	var image Image
	if err := db.First(&image, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Serve the image file.
	c.File("uploads/" + id)
}
