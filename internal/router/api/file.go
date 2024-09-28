package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FileUploadRequest represents the request body for uploading a file
// @Description FileUploadRequest represents the request body for uploading a file
type FileUploadRequest struct {
	// File is the file to be uploaded
	// @Description File is the file to be uploaded
	File *multipart.FileHeader `form:"file" binding:"required"`
	// SaveDir is the directory where the file will be saved
	// @Description SaveDir is the directory where the file will be saved
	SaveDir string `form:"directory" binding:"required"`
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to the specified directory, please just use name of dir e.g. "dirName"
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param directory formData string true "Directory to save the file"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /upload [post]
func UploadFile(c *gin.Context) {
	var req FileUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		logrus.Warningf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "Invalid request", "error": err.Error()},
		)
	}

	file, err := c.FormFile("file")
	if err != nil {
		logrus.Warningf("Failed to get file: %v", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "Failed to get file", "error": err},
		)
		return
	}

	var fileName string

	if strings.HasSuffix(file.Filename, ".yml") {
		fileName = "docker-compose.yml"
	} else if strings.HasSuffix(file.Filename, ".env") {
		fileName = ".env"
	} else {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "File extension must be .env or .yml"},
		)
		return
	}

	filePath := fmt.Sprintf("uploads/%s/%s", req.SaveDir, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logrus.Warningf("Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully."})
}
