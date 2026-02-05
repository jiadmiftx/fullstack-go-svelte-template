package unit_settings_controller

import (
	"net/http"
	"time"

	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/gin_utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GinUnitSettingsController interface {
	GetSettings(c *gin.Context)
	UpdateSettings(c *gin.Context)
}

type controller struct {
	db *gorm.DB
}

func NewUnitSettingsController(db *gorm.DB) GinUnitSettingsController {
	return &controller{db: db}
}

// GetSettings returns unit settings or creates default if not exists
func (ctrl *controller) GetSettings(c *gin.Context) {
	unitIdStr := c.Param("id")
	unitId, err := uuid.Parse(unitIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid unit id"})
		return
	}

	var settings schemas.UnitSettings
	if err := ctrl.db.Where("unit_id = ?", unitId).First(&settings).Error; err != nil {
		// Create default settings
		settings = schemas.UnitSettings{
			Id:               uuid.New(),
			UnitId:           unitId,
			PeriodDuration:   40,
			StartTime:        "07:00",
			TotalPeriods:     9,
			BreakAfterPeriod: 3,
			BreakDuration:    15,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := ctrl.db.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "success",
		Data: map[string]interface{}{
			"id":                 settings.Id.String(),
			"unit_id":            settings.UnitId.String(),
			"period_duration":    settings.PeriodDuration,
			"start_time":         settings.StartTime,
			"total_periods":      settings.TotalPeriods,
			"break_after_period": settings.BreakAfterPeriod,
			"break_duration":     settings.BreakDuration,
			"academic_year":      settings.AcademicYear,
			"current_semester":   settings.CurrentSemester,
			"semester_1_start":   settings.Semester1Start,
			"semester_1_end":     settings.Semester1End,
			"semester_2_start":   settings.Semester2Start,
			"semester_2_end":     settings.Semester2End,
		},
	})
}

type UpdateSettingsRequest struct {
	PeriodDuration   *int    `json:"period_duration"`
	StartTime        *string `json:"start_time"`
	TotalPeriods     *int    `json:"total_periods"`
	BreakAfterPeriod *int    `json:"break_after_period"`
	BreakDuration    *int    `json:"break_duration"`
	AcademicYear     *string `json:"academic_year"`
	CurrentSemester  *int    `json:"current_semester"`
	Semester1Start   *string `json:"semester_1_start"`
	Semester1End     *string `json:"semester_1_end"`
	Semester2Start   *string `json:"semester_2_start"`
	Semester2End     *string `json:"semester_2_end"`
}

// UpdateSettings updates unit settings
func (ctrl *controller) UpdateSettings(c *gin.Context) {
	unitIdStr := c.Param("id")
	unitId, err := uuid.Parse(unitIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "invalid unit id"})
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var settings schemas.UnitSettings
	if err := ctrl.db.Where("unit_id = ?", unitId).First(&settings).Error; err != nil {
		settings = schemas.UnitSettings{
			Id:     uuid.New(),
			UnitId: unitId,
		}
	}

	if req.PeriodDuration != nil {
		settings.PeriodDuration = *req.PeriodDuration
	}
	if req.StartTime != nil {
		settings.StartTime = *req.StartTime
	}
	if req.TotalPeriods != nil {
		settings.TotalPeriods = *req.TotalPeriods
	}
	if req.BreakAfterPeriod != nil {
		settings.BreakAfterPeriod = *req.BreakAfterPeriod
	}
	if req.BreakDuration != nil {
		settings.BreakDuration = *req.BreakDuration
	}
	if req.AcademicYear != nil {
		settings.AcademicYear = *req.AcademicYear
	}
	if req.CurrentSemester != nil {
		settings.CurrentSemester = *req.CurrentSemester
	}
	if req.Semester1Start != nil {
		if t, err := time.Parse("2006-01-02", *req.Semester1Start); err == nil {
			settings.Semester1Start = &t
		}
	}
	if req.Semester1End != nil {
		if t, err := time.Parse("2006-01-02", *req.Semester1End); err == nil {
			settings.Semester1End = &t
		}
	}
	if req.Semester2Start != nil {
		if t, err := time.Parse("2006-01-02", *req.Semester2Start); err == nil {
			settings.Semester2Start = &t
		}
	}
	if req.Semester2End != nil {
		if t, err := time.Parse("2006-01-02", *req.Semester2End); err == nil {
			settings.Semester2End = &t
		}
	}

	if err := ctrl.db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "settings updated",
		Data:    map[string]string{"id": settings.Id.String()},
	})
}
