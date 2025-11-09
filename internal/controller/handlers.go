package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Oleska1601/WBURLShortener/internal/apperrors"
	"github.com/Oleska1601/WBURLShortener/internal/controller/dto"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
)

// CreateShortURLHandler godoc
// @Summary create short url
// @Description create short url with provided params
// @Tags url
// @Accept json
// @Produce json
// @Param notification body dto.CreateShortURLRequest true "create short url request"
// @Success	201	{object} CreateShortURLResponse
// @Failure	400	{object} map[string]string "impossible to create short url"
// @Failure 409	{object} map[string]string "provided short url already exists"
// @Failure	500	{object} map[string]string "failed to create notification"
// @Router /api/shorten [post]
func (s *Server) CreateShortURLHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "CreateShortURLHandler c.ShouldBindJSON").
			Msg("impossible to create short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("impossible to create short url")})
		return
	}

	modelURL, err := req.ToModel()
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "CreateShortURLHandler req.ToModel").
			Msg("impossible to create short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := s.usecase.CreateURL(ctx, modelURL)
	if err != nil {
		if errors.Is(err, apperrors.AlreadyExistsError) {
			zlog.Logger.Error().
				Err(err).
				Int("status", http.StatusConflict).
				Str("path", "CreateShortURLHandler s.usecase.CreateShortURL").
				Msg("conflict to create short url")

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "CreateShortURLHandler s.usecase.CreateShortURL").
			Msg("failed to create short url")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create short url"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "CreateShortURLHandler").
		Str("short_url", url.ShortURL).
		Msg("create short url successful")

	c.JSON(http.StatusCreated, toCreateShortURLResponse(url))
}

// RedirectShortURLHandler godoc
// @Summary redirect by short url
// @Description redirect to original url by provided short url
// @Tags url
// @Accept json
// @Produce json
// @Param short_url path string true "short url"
// @Success	302 "redirect by short url"
// @Failure	400	{object} map[string]string "invalid short url"
// @Failure	404	{object} map[string]string "provided short url does not exist"
// @Failure	500	{object} map[string]string "failed to redirect by provided short url"
// @Router /api/s/{short_url} [get]
func (s *Server) RedirectShortURLHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userAgent := c.Request.Header.Get("User-Agent")
	if userAgent == "" {
		zlog.Logger.Warn().
			Str("path", "RedirectShortURLHandler c.Request.Header.Get").
			Msg("unknown User-Agent")
		userAgent = "unknown:User-Agent"
	}

	ip := c.ClientIP()
	if ip == "" {
		zlog.Logger.Warn().
			Str("path", "RedirectShortURLHandler c.ClientIP").
			Msg("unknown IP")
		ip = "unknown:IP"
	}

	shortURL := strings.TrimSpace(c.Param("short_url"))
	if shortURL == "" {
		zlog.Logger.Error().
			Int("status", http.StatusBadRequest).
			Str("path", "GetNotificationStatusHandler").
			Msg("invalid short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short url"})
		return
	}

	fmt.Println(shortURL)

	originalURL, err := s.usecase.GetURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, apperrors.NotFoundError) {
			zlog.Logger.Error().
				Err(err).
				Int("status", http.StatusNotFound).
				Str("path", "CreateShortURLHandler s.usecase.CreateShortURL").
				Msg("cannot get short url")

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "CreateShortURLHandler s.usecase.CreateShortURL").
			Msg("failed to get short url")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get short url"})
		return
	}
	zlog.Logger.Info().
		Int("status", http.StatusFound).
		Str("path", "RedirectShortURLHandler").
		Str("short_url", shortURL).
		Msg("redirect by short_url successful")

	// create analytics
	analytics := dto.CreateAnalytics(shortURL, userAgent, ip)
	id, err := s.usecase.CreateAnalytics(ctx, analytics)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("path", "s.usecase.CreateAnalytics").
			Str("short_url", analytics.ShortURL).
			Int64("id", id).
			Msg("failed to create analytics")
	} else {
		zlog.Logger.Info().
			Str("path", "RedirectShortURLHandler CreateAnalytics").
			Str("short_url", analytics.ShortURL).
			Int64("id", id).
			Msg("create analytics successful")
	}

	c.Redirect(http.StatusFound, originalURL)
}

// GetAnalyticsHandler godoc
// @Summary get analytics by short url
// @Description get analytics by provided short url
// @Tags notify
// @Accept json
// @Produce json
// @Param short_url path string true "short url"
// @Success 200	{object} GetAnalyticsResponse
// @Failure 400	{object} map[string]string "invalid short url"
// @Failure	500	{object} map[string]string "failed to get analytics"
// @Router /api/analytics/{short_url} [get]
func (s *Server) GetAnalyticsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	shortURL := strings.TrimSpace(c.Param("short_url"))
	if shortURL == "" {
		zlog.Logger.Error().
			Int("status", http.StatusBadRequest).
			Str("path", "GetAnalyticsHandler").
			Msg("invalid short url")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short url"})
		return
	}

	analytics, err := s.usecase.GetAnalytics(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "GetAnalyticsHandler s.usecase.GetAnalytics").
			Msg("failed to get analytics")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get analytics"})
		return
	}

	fmt.Println(analytics)

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "GetAnalyticsHandler").
		Str("short_url", shortURL).
		Msg("get analytics by short url successful")

	c.JSON(http.StatusOK, toAnalyticsResponse(analytics))
}
