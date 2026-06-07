package middleware

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func ValidateClassMajorRelation(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		classIDStr := c.Query("class_id")
		majorIDStr := c.Query("major_id")
		ayIDStr := c.Query("academic_year_id")
		if ayIDStr == "" {
			ayIDStr = c.Query("academic_year")
		}
		entryYearStr := c.Query("entry_year")

		var classID, majorID, academicYearID uint

		if classIDStr != "" {
			if cid, err := strconv.ParseUint(classIDStr, 10, 32); err == nil {
				classID = uint(cid)
			}
		}
		if majorIDStr != "" {
			if mid, err := strconv.ParseUint(majorIDStr, 10, 32); err == nil {
				majorID = uint(mid)
			}
		}

		// Resolve academicYearID
		if ayIDStr != "" {
			if ayid, err := strconv.ParseUint(ayIDStr, 10, 32); err == nil {
				academicYearID = uint(ayid)
			}
		} else if entryYearStr != "" {
			if ey, err := strconv.Atoi(entryYearStr); err == nil && ey > 0 {
				var ay domain.AcademicYear
				err := db.NewSelect().
					Model(&ay).
					Where("year = ?", ey).
					Limit(1).
					Scan(c.Request.Context())
				if err == nil {
					academicYearID = ay.ID
				} else {
					helper.ErrorResponse(c, http.StatusBadRequest, "Angkatan tidak ditemukan")
					c.Abort()
					return
				}
			}
		}

		// 1. Validasi Class - Major
		if classID > 0 && majorID > 0 {
			var cls domain.Class
			err := db.NewSelect().
				Model(&cls).
				Where("id = ?", classID).
				Scan(c.Request.Context())
			if err != nil {
				helper.ErrorResponse(c, http.StatusBadRequest, "Kelas yang dipilih tidak ditemukan")
				c.Abort()
				return
			}
			if cls.MajorID == nil || *cls.MajorID != majorID {
				helper.ErrorResponse(c, http.StatusBadRequest, "Kelas yang dipilih tidak sesuai dengan Jurusan")
				c.Abort()
				return
			}
		}

		// 2. Validasi AcademicYear - Major
		if academicYearID > 0 && majorID > 0 {
			var aym domain.AcademicYearMajor
			exists, err := db.NewSelect().
				Model(&aym).
				Where("academic_year_id = ? AND major_id = ?", academicYearID, majorID).
				Exists(c.Request.Context())
			if err != nil {
				helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memverifikasi relasi angkatan dan jurusan")
				c.Abort()
				return
			}
			if !exists {
				helper.ErrorResponse(c, http.StatusBadRequest, "Jurusan yang dipilih tidak terdaftar pada Angkatan ini")
				c.Abort()
				return
			}
		}

		// 3. Validasi AcademicYear - Class
		if academicYearID > 0 && classID > 0 {
			var ayc domain.AcademicYearClass
			exists, err := db.NewSelect().
				Model(&ayc).
				Where("academic_year_id = ? AND class_id = ?", academicYearID, classID).
				Exists(c.Request.Context())
			if err != nil {
				helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memverifikasi relasi angkatan dan kelas")
				c.Abort()
				return
			}
			if !exists {
				helper.ErrorResponse(c, http.StatusBadRequest, "Kelas yang dipilih tidak terdaftar pada Angkatan ini")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
