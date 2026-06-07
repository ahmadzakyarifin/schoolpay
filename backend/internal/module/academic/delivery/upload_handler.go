package delivery

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadStudentPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "file foto tidak ditemukan")
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		helper.ErrorResponse(c, http.StatusBadRequest, "format file harus jpg, jpeg, atau png")
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("public/uploads/students", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "gagal menyimpan file")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "foto berhasil diupload", gin.H{
		"path": "/uploads/students/" + filename,
	})
}
