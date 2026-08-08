package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/businessdate"
	httpx "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var chinaStandardTime = time.FixedZone("Asia/Shanghai", 8*60*60)

type DailyTurnoverController struct {
	service *service.DailyTurnoverService
}

func NewDailyTurnoverController(service *service.DailyTurnoverService) *DailyTurnoverController {
	return &DailyTurnoverController{service: service}
}

func (c *DailyTurnoverController) List(ctx *gin.Context) {
	businessDate, err := resolveDailyTurnoverBusinessDate(ctx.Query("business_date"), time.Now())
	if err != nil {
		internalControllerError(ctx, http.StatusBadRequest, "business_date must use YYYY-MM-DD format")
		return
	}

	reports, err := c.service.List(ctx.Request.Context(), businessDate)
	if err != nil {
		logging.LogError("daily turnover report query failed", zap.Error(err))
		internalControllerError(ctx, http.StatusInternalServerError, "internal server error")
		return
	}
	httpx.Success(ctx, reports)
}

func resolveDailyTurnoverBusinessDate(raw string, now time.Time) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return businessdate.Date(now.In(chinaStandardTime)).AddDate(0, 0, -1).Format("2006-01-02"), nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", raw, chinaStandardTime)
	if err != nil {
		return "", err
	}
	return parsed.Format("2006-01-02"), nil
}

func internalControllerError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, httpx.Response{
		Code:    status,
		Message: message,
	})
}
