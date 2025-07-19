package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/JonathanVera18/ecommerce-api/internal/utils"
)

type FileUploadHandler struct {
	uploadDir string
}

func NewFileUploadHandler(uploadDir string) *FileUploadHandler {
	return &FileUploadHandler{uploadDir: uploadDir}
}

// UploadFile handles file uploads (images, documents, etc.)
func (h *FileUploadHandler) UploadFile(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "No files uploaded")
	}

	uploadedFiles := []map[string]string{}
	maxFileSize := int64(10 << 20) // 10MB
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"application/pdf": true,
		"text/plain": true,
	}

	for _, file := range files {
		// Check file size
		if file.Size > maxFileSize {
			return utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("File %s is too large (max 10MB)", file.Filename))
		}

		// Check file type
		src, err := file.Open()
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to open uploaded file")
		}
		defer src.Close()

		// Read first 512 bytes to detect content type
		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		}

		contentType := http.DetectContentType(buffer)
		if !allowedTypes[contentType] {
			return utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("File type %s not allowed", contentType))
		}

		// Reset file reader
		src.Seek(0, io.SeekStart)

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
		
		// Create directory if it doesn't exist
		userDir := filepath.Join(h.uploadDir, fmt.Sprintf("user_%d", userID))
		if err := os.MkdirAll(userDir, 0755); err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create upload directory")
		}

		// Save file
		dst, err := os.Create(filepath.Join(userDir, filename))
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create destination file")
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		}

		// Add to uploaded files list
		uploadedFiles = append(uploadedFiles, map[string]string{
			"filename":     filename,
			"original_name": file.Filename,
			"content_type": contentType,
			"size":        fmt.Sprintf("%d", file.Size),
			"url":         fmt.Sprintf("/uploads/user_%d/%s", userID, filename),
		})
	}

	return utils.CreatedResponse(c, "Files uploaded successfully", uploadedFiles)
}

// GetUserFiles retrieves files uploaded by a user
func (h *FileUploadHandler) GetUserFiles(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	userDir := filepath.Join(h.uploadDir, fmt.Sprintf("user_%d", userID))
	files := []map[string]string{}

	// Check if directory exists
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return utils.SuccessResponse(c, "User files retrieved successfully", files)
	}

	// Read directory contents
	entries, err := os.ReadDir(userDir)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read user directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, map[string]string{
				"filename":    entry.Name(),
				"size":       fmt.Sprintf("%d", info.Size()),
				"modified":   info.ModTime().Format(time.RFC3339),
				"url":        fmt.Sprintf("/uploads/user_%d/%s", userID, entry.Name()),
			})
		}
	}

	return utils.SuccessResponse(c, "User files retrieved successfully", files)
}

// DeleteFile deletes a user's uploaded file
func (h *FileUploadHandler) DeleteFile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	filename := c.Param("filename")

	if filename == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Filename is required")
	}

	// Sanitize filename to prevent directory traversal
	filename = filepath.Base(filename)
	if strings.Contains(filename, "..") {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filename")
	}

	filePath := filepath.Join(h.uploadDir, fmt.Sprintf("user_%d", userID), filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return utils.ErrorResponse(c, http.StatusNotFound, "File not found")
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete file")
	}

	return utils.SuccessResponse(c, "File deleted successfully", nil)
}

// ServeFile serves uploaded files
func (h *FileUploadHandler) ServeFile(c echo.Context) error {
	userID := c.Param("userId")
	filename := c.Param("filename")

	if userID == "" || filename == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "User ID and filename are required")
	}

	// Sanitize filename to prevent directory traversal
	filename = filepath.Base(filename)
	if strings.Contains(filename, "..") {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filename")
	}

	filePath := filepath.Join(h.uploadDir, fmt.Sprintf("user_%s", userID), filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return utils.ErrorResponse(c, http.StatusNotFound, "File not found")
	}

	return c.File(filePath)
}
