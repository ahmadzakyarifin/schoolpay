package delivery

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FinanceReportHandler struct {
	s usecase.FinanceReportService
}

func NewFinanceReportHandler(s usecase.FinanceReportService) *FinanceReportHandler {
	return &FinanceReportHandler{s: s}
}

func (h *FinanceReportHandler) GetArrears(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	academicYear, _ := strconv.Atoi(c.Query("academic_year_id"))
	if academicYear == 0 {
		academicYear, _ = strconv.Atoi(c.Query("academic_year"))
	}
	classID, _ := strconv.Atoi(c.Query("class_id"))
	majorID, _ := strconv.Atoi(c.Query("major_id"))
	billTypeID, _ := strconv.Atoi(c.Query("bill_type_id"))
	search := c.Query("search")
	period := c.DefaultQuery("period", "all")

	now := time.Now()
	var start, end *time.Time

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		sStart := time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		eEnd := time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
		start = &sStart
		end = &eEnd
	}

	refTime := now
	refDateStr := c.Query("ref_date")
	if refDateStr != "" {
		switch period {
		case "monthly":
			if t, err := time.Parse("2006-01", refDateStr); err == nil {
				refTime = t
			}
		case "daily":
			if t, err := time.Parse("2006-01-02", refDateStr); err == nil {
				refTime = t
			}
		case "yearly":
			if t, err := time.Parse("2006", refDateStr); err == nil {
				refTime = t
			}
		}
	}

	if period != "custom" && period != "all" {
		switch period {
		case "daily":
			sStart := time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			eEnd := time.Date(sStart.Year(), sStart.Month(), sStart.Day(), 23, 59, 59, 0, sStart.Location())
			start = &sStart
			end = &eEnd
		case "yearly":
			sStart := time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			eEnd := time.Date(sStart.Year(), 12, 31, 23, 59, 59, 0, sStart.Location())
			start = &sStart
			end = &eEnd
		default: // monthly
			sStart := time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			eEnd := sStart.AddDate(0, 1, 0).Add(-time.Second)
			start = &sStart
			end = &eEnd
		}
	}

	list, total, err := h.s.GetArrears(c.Request.Context(), page, limit, academicYear, uint(classID), uint(majorID), uint(billTypeID), search, start, end)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *FinanceReportHandler) ExportTrend(c *gin.Context) {
	interval := c.DefaultQuery("interval", "daily")
	academicYear, _ := strconv.Atoi(c.Query("academic_year"))
	classID, _ := strconv.Atoi(c.Query("class_id"))
	majorID, _ := strconv.Atoi(c.Query("major_id"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var start, end *time.Time
	if startDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		start = &s
	}
	if endDateStr != "" {
		e, _ := time.Parse("2006-01-02", endDateStr)
		end = &e
	}

	data, err := h.s.ExportTrendExcel(c.Request.Context(), start, end, interval, academicYear, uint(classID), uint(majorID))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=payment_trend.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
