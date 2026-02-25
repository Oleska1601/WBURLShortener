package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Oleska1601/WBURLShortener/internal/controller/api/v1/request"
	"github.com/Oleska1601/WBURLShortener/internal/controller/api/v1/response"
	"github.com/Oleska1601/WBURLShortener/internal/errs"
	"github.com/Oleska1601/WBURLShortener/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
)

const (
	createShortURLURI   = "/shorten"
	redirectShortURLURI = "/s/:short_url"
	getAnalyticsURI     = "/analytics/:short_url"
)

func (v1 *APIV1) registerHandlers(group *gin.RouterGroup) {
	group.POST(createShortURLURI, v1.createShortURL)
	group.GET(redirectShortURLURI, v1.redirectShortURL)
	group.GET(getAnalyticsURI, v1.getAnalytics)
}

// @Summary create short url
// @Description create short url with provided params
// @Tags URL API
// @Accept json
// @Produce json
// @Param request body request.CreateShortURLRequest true "create short url request"
// @Success	201	{object} response.CreateShortURLResponse
// @Failure	400	{object} map[string]string "validate error"
// @Failure 409	{object} map[string]string "conflict error"
// @Failure	500	{object} map[string]string "server error"
// @Router /api/v1/shorten [post]
func (v1 *APIV1) createShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "createShortURL c.ShouldBindJSON").
			Msg("impossible to create short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("impossible to create short url")})
		return
	}

	if err := validate.StructCtx(ctx, req); err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "createShortURL validate.StructCtx").
			Msg("impossible to create short url: validator error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to create short url: validator error"})
		return
	}

	modelURL, err := req.ToModel()
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "createShortURL req.ToModel").
			Msg("impossible to create short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("impossible to create short url: %v", err.Error())})
		return
	}

	url, err := v1.service.CreateURL(ctx, modelURL)
	if err != nil {
		if errors.Is(err, errs.AlreadyExistsError) || errors.Is(err, errs.ConflictError) {
			zlog.Logger.Error().
				Err(err).
				Int("status", http.StatusConflict).
				Str("path", "createShortURL v1.service.CreateURL").
				Msg("failed to create short url: conflict error")
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("failed to create short url: %v", err.Error())})
			return
		}

		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "createShortURL v1.service.CreateURL").
			Msg("failed to create short url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create short url"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusCreated).
		Str("path", "createShortURL").
		Str("short_url", url.ShortURL).
		Msg("create short url successful")
	c.JSON(http.StatusCreated, response.ToCreateShortURLResponse(url))
}

// @Summary redirect by short url
// @Description redirect to original url by provided short url
// @Tags URL API
// @Produce json,text/plain,html
// @Param short_url path string true "short url"
// @Success	302 "redirect by short url"
// @Header 302 {string} Location "redirect URL"
// @Failure	400	{object} map[string]string "invalid short url"
// @Failure	404	{object} map[string]string "provided short url does not exist"
// @Failure	500	{object} map[string]string "server error"
// @Router /api/v1/s/{short_url} [get]
func (v1 *APIV1) redirectShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	userAgent := c.Request.Header.Get("User-Agent")
	if userAgent == "" {
		zlog.Logger.Warn().
			Str("path", "redirectShortURL c.Request.Header.Get").
			Msg("unknown User-Agent")
		userAgent = "unknown:User-Agent"
	}

	ip := c.ClientIP()
	if ip == "" {
		zlog.Logger.Warn().
			Str("path", "redirectShortURL c.ClientIP").
			Msg("unknown IP")
		ip = "unknown:IP"
	}

	shortURL := strings.TrimSpace(c.Param("short_url"))
	if shortURL == "" {
		zlog.Logger.Error().
			Int("status", http.StatusBadRequest).
			Str("path", "redirectShortURL").
			Msg("impossible to redirect by short url: invalid short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to redirect by short url: invalid short url"})
		return
	}

	originalURL, err := v1.service.GetURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, errs.NotFoundError) {
			zlog.Logger.Error().
				Err(err).
				Int("status", http.StatusNotFound).
				Str("path", "redirectShortURL v1.service.GetURL").
				Str("short_url", shortURL).
				Msg("failed to redirect by short url: not found error")
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to redirect by short url: %v", err.Error())})
			return
		}

		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "redirectShortURL v1.service.GetURL").
			Str("short_url", shortURL).
			Msg("failed to redirect by short url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to redirect by short url"})
		return
	}
	zlog.Logger.Info().
		Int("status", http.StatusFound).
		Str("path", "redirectShortURL").
		Str("short_url", shortURL).
		Msg("redirect by short_url successful")

	analytics := &models.Analytics{
		ShortURL:  shortURL,
		UserAgent: userAgent,
		IP:        ip,
	}
	id, err := v1.service.CreateAnalytics(ctx, analytics)
	if err != nil {
		// ошибка пользователю не возвращается, тк это создание аналитики - не главная задача endpoint-а
		zlog.Logger.Error().
			Err(err).
			Str("path", "redirectShortURL v1.service.CreateAnalytics").
			Str("short_url", analytics.ShortURL).
			Msg("failed to create analytics")
	} else {
		zlog.Logger.Info().
			Str("path", "redirectShortURL v1.service.CreateAnalytics").
			Str("short_url", analytics.ShortURL).
			Int("id", id).
			Msg("create analytics successful")
	}

	c.Redirect(http.StatusFound, originalURL)
}

// @Summary get analytics by short url
// @Description get analytics by provided short url
// @Tags ANALYTICS API
// @Produce json
// @Param short_url path string true "short url"
// @Success 200	{object} models.AnAgregation
// @Failure 400	{object} map[string]string "invalid short url"
// @Failure	500	{object} map[string]string "server error"
// @Router /api/v1/analytics/{short_url} [get]
func (v1 *APIV1) getAnalytics(c *gin.Context) {
	ctx := c.Request.Context()
	shortURL := strings.TrimSpace(c.Param("short_url"))
	if shortURL == "" {
		zlog.Logger.Error().
			Int("status", http.StatusBadRequest).
			Str("path", "getAnalytics").
			Msg("impossible to get analytics: invalid short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to get analytics: invalid short url"})
		return
	}

	analytics, err := v1.service.GetAnalytics(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "getAnalytics v1.service.GetAnalytics").
			Str("short_url", shortURL).
			Msg("failed to get analytics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get analytics"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "getAnalytics").
		Str("short_url", shortURL).
		Msg("get analytics by short url successful")
	c.JSON(http.StatusOK, analytics)
}
