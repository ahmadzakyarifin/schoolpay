package delivery

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashSvc financeusecase.DashboardService
}

func NewDashboardHandler(
	dashSvc financeusecase.DashboardService,
) *DashboardHandler {
	return &DashboardHandler{dashSvc}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()
	period := c.DefaultQuery("period", "all")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	now := time.Now()
	var start, end, prevStart, prevEnd time.Time

	// 1. Resolve Filter Params
	var filterEntryYear int
	if c.Query("academic_year_id") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("academic_year_id"))
	} else if c.Query("entry_year") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("entry_year"))
	}
	classID, _ := strconv.ParseUint(c.Query("class_id"), 10, 32)
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)

	// 2. Handle Periods
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		start = time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		end = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
		diff := end.Sub(start)
		prevEnd = start
		prevStart = start.Add(-diff)
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
			start = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
			prevStart = start.Add(-24 * time.Hour)
			prevEnd = time.Date(prevStart.Year(), prevStart.Month(), prevStart.Day(), 23, 59, 59, 0, prevStart.Location())
		case "yearly":
			start = time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), 12, 31, 23, 59, 59, 0, start.Location())
			prevStart = start.AddDate(-1, 0, 0)
			prevEnd = time.Date(prevStart.Year(), 12, 31, 23, 59, 59, 0, prevStart.Location())
		default: // monthly
			start = time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Second)
			prevStart = start.AddDate(0, -1, 0)
			prevEnd = prevStart.AddDate(0, 1, 0).Add(-time.Second)
		}
	}

	var pStart, pEnd, pPrevStart, pPrevEnd *time.Time
	if period != "all" && !(period == "custom" && (startDateStr == "" || endDateStr == "")) {
		pStart = &start
		pEnd = &end
		pPrevStart = &prevStart
		pPrevEnd = &prevEnd
	}

	billTypeID, _ := strconv.ParseUint(c.Query("bill_type_id"), 10, 32)

	statsData, err := h.dashSvc.GetStats(ctx, period, pStart, pEnd, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID), page, limit, uint(billTypeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   statsData,
	})
}

func (h *DashboardHandler) GetCommunicationDetails(c *gin.Context) {
	ctx := c.Request.Context()
	status := strings.TrimSpace(strings.ToLower(c.Query("status")))
	channel := c.DefaultQuery("channel", "whatsapp")
	period := c.DefaultQuery("period", "all")

	// 1. Resolve Filter Params
	var filterEntryYear int
	if c.Query("academic_year_id") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("academic_year_id"))
	} else if c.Query("entry_year") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("entry_year"))
	}
	classID, _ := strconv.ParseUint(c.Query("class_id"), 10, 32)
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)

	// 2. Handle Periods
	var start, end time.Time
	now := time.Now()
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		start = time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		end = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
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
			start = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		case "yearly":
			start = time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), 12, 31, 23, 59, 59, 0, start.Location())
		default: // monthly
			start = time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Second)
		}
	}

	var pStart, pEnd *time.Time
	if period != "all" && !(period == "custom" && (startDateStr == "" || endDateStr == "")) {
		pStart = &start
		pEnd = &end
	}

	list, err := h.dashSvc.GetCommunicationDetails(ctx, status, channel, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *DashboardHandler) ExportGlobalReport(c *gin.Context) {
	ctx := c.Request.Context()
	tab := c.Query("tab")
	search := c.Query("search")
	period := c.DefaultQuery("period", "all")
	format := c.Query("format")

	now := time.Now()
	var start, end time.Time

	var filterEntryYear int
	if c.Query("academic_year_id") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("academic_year_id"))
	} else if c.Query("entry_year") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("entry_year"))
	}
	classID, _ := strconv.ParseUint(c.Query("class_id"), 10, 32)
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)
	billTypeID, _ := strconv.ParseUint(c.Query("bill_type_id"), 10, 32)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		start = time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		end = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
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
			start = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		case "yearly":
			start = time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), 12, 31, 23, 59, 59, 0, start.Location())
		default: // monthly
			start = time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Second)
		}
	}

	var pStart, pEnd *time.Time
	if period != "all" && !(period == "custom" && (startDateStr == "" || endDateStr == "")) {
		pStart = &start
		pEnd = &end
	}

	data, err := h.dashSvc.ExportGlobalReport(ctx, format, tab, search, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), uint(billTypeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if format == "pdf" {
		c.Header("Content-Disposition", "attachment; filename=Laporan_Keuangan.pdf")
		c.Header("Content-Type", "application/pdf")
		c.Data(http.StatusOK, "application/pdf", data)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Laporan_Keuangan.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
