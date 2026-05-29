package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	s   usecase.StudentService
	cfg *config.Config
}

func NewStudentHandler(s usecase.StudentService, cfg *config.Config) *StudentHandler {
	return &StudentHandler{s: s, cfg: cfg}
}

func (h *StudentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filter := c.Query("filter")
	status := c.Query("status")
	entryYear, _ := strconv.Atoi(c.Query("entry_year"))
	classID, _ := strconv.Atoi(c.Query("class_id"))
	majorID, _ := strconv.Atoi(c.Query("major_id"))
	sort := c.Query("sort")

	list, total, err := h.s.GetPaginated(c.Request.Context(), page, limit, search, filter, status, entryYear, uint(classID), uint(majorID), sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengambil data siswa: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{
		"data":  list,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *StudentHandler) GetParents(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	list, err := h.s.GetParents(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengambil data orang tua")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *StudentHandler) GetFilters(c *gin.Context) {
	data, err := h.s.GetAcademicFilters(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengambil data filter: "+err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", data)
}

func (h *StudentHandler) Create(c *gin.Context) {
	var student domain.Student
	allErrors := make(map[string][]string)

	// 1. Initial Binding & Basic Validation
	if err := c.ShouldBind(&student); err != nil {
		allErrors = utils.GetValidationErrors(err)
		// We don't return early here to allow business duplicate checks to run
	}

	// 2. Handle Image Upload (only if no critical binding error or we can still proceed)
	file, _ := c.FormFile("image_path")
	if file != nil {
		filename := fmt.Sprintf("student_%s_%s", student.NISN, file.Filename)
		dst := fmt.Sprintf("public/uploads/students/%s", filename)
		if err := h.ensureDir("public/uploads/students"); err == nil {
			if err := c.SaveUploadedFile(file, dst); err == nil {
				imgPath := "uploads/students/" + filename
				student.ImagePath = &imgPath
			}
		}
	}

	// 3. Service Level Validation & Creation
	// We call Create even if allErrors has content, because Create will perform duplicate checks
	if err := h.s.Create(c.Request.Context(), &student); err != nil {
		if bmErr, ok := err.(*utils.BusinessMultiError); ok {
			for k, v := range bmErr.Errors {
				allErrors[k] = append(allErrors[k], v...)
			}
		} else if bErr, ok := err.(*utils.BusinessError); ok {
			allErrors[bErr.Field] = append(allErrors[bErr.Field], bErr.Message)
		} else {
			// If it's a real unexpected error and we have no validation errors yet, return 500
			if len(allErrors) == 0 {
				utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// 4. Final Error Report
	if len(allErrors) > 0 {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", allErrors)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "data siswa berhasil dibuat", student)
}

func (h *StudentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	var student domain.Student
	allErrors := make(map[string][]string)

	// 1. Initial Binding & Basic Validation
	if err := c.ShouldBind(&student); err != nil {
		allErrors = utils.GetValidationErrors(err)
		// We don't return early here
	}

	// 2. Handle Image Upload
	file, _ := c.FormFile("image_path")
	if file != nil {
		filename := fmt.Sprintf("student_%s_%s", student.NISN, file.Filename)
		dst := fmt.Sprintf("public/uploads/students/%s", filename)
		if err := h.ensureDir("public/uploads/students"); err == nil {
			if err := c.SaveUploadedFile(file, dst); err == nil {
				imgPath := "uploads/students/" + filename
				student.ImagePath = &imgPath
			}
		}
	}

	id, _ := strconv.Atoi(idStr)
	student.ID = uint(id)

	// 3. Service Level Validation & Update
	if err := h.s.Update(c.Request.Context(), &student); err != nil {
		if bmErr, ok := err.(*utils.BusinessMultiError); ok {
			for k, v := range bmErr.Errors {
				allErrors[k] = append(allErrors[k], v...)
			}
		} else if bErr, ok := err.(*utils.BusinessError); ok {
			allErrors[bErr.Field] = append(allErrors[bErr.Field], bErr.Message)
		} else {
			if len(allErrors) == 0 {
				utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// 4. Final Error Report
	if len(allErrors) > 0 {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", allErrors)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "data siswa berhasil diperbarui", student)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	student, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "data siswa tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil", student)
}

func (h *StudentHandler) ToggleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengubah status siswa")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "status siswa berhasil diubah", nil)
}

func (h *StudentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "data siswa berhasil dihapus", nil)
}

func (h *StudentHandler) GetClassHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	list, err := h.s.GetClassHistory(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengambil riwayat kelas: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *StudentHandler) GetMyStudents(c *gin.Context) {
	userID, _ := c.Get("user_id")
	list, err := h.s.GetStudentsByParentID(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengambil data anak: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *StudentHandler) GetByParentID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	list, err := h.s.GetStudentsByParentID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "gagal mengambil data anak: "+err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *StudentHandler) Export(c *gin.Context) {
	ctx := c.Request.Context()
	search := c.Query("search")
	filter := c.Query("filter")
	status := c.Query("status")
	entryYear, _ := strconv.Atoi(c.Query("entry_year"))
	classID, _ := strconv.Atoi(c.Query("class_id"))
	majorID, _ := strconv.Atoi(c.Query("major_id"))

	data, err := h.s.ExportExcel(ctx, search, filter, status, entryYear, uint(classID), uint(majorID))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Daftar_Siswa.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *StudentHandler) BulkGraduate(c *gin.Context) {
	var body struct {
		ClassID    uint   `json:"class_id"`
		StudentIDs []uint `json:"student_ids"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "data tidak valid")
		return
	}

	if err := h.s.BulkGraduate(c.Request.Context(), body.ClassID, body.StudentIDs); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Kelulusan masal berhasil diproses", nil)
}

func (h *StudentHandler) BulkDelete(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.Ids); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "data siswa terpilih berhasil dinonaktifkan", nil)
}

func (h *StudentHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "data siswa berhasil dipulihkan", nil)
}

func (h *StudentHandler) BulkRestore(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.Ids); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d data siswa berhasil dipulihkan", len(req.Ids)), nil)
}

func (h *StudentHandler) BulkPromote(c *gin.Context) {
	var body struct {
		SourceClassID uint   `json:"source_class_id"`
		TargetClassID uint   `json:"target_class_id"`
		StudentIDs    []uint `json:"student_ids"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "data tidak valid")
		return
	}

	if err := h.s.BulkPromote(c.Request.Context(), body.SourceClassID, body.TargetClassID, body.StudentIDs); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perpindahan kelas masal berhasil diproses", nil)
}

func (h *StudentHandler) ensureDir(path string) error {
	return utils.EnsureDir(path)
}

func (h *StudentHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *StudentHandler) CheckUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeIDStr := c.DefaultQuery("exclude_id", "0")
	excludeID, _ := strconv.Atoi(excludeIDStr)

	isUnique, err := h.s.CheckUnique(c.Request.Context(), field, value, uint(excludeID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{
		"is_unique": isUnique,
	})
}

