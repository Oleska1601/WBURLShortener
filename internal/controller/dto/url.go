package dto

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/Oleska1601/WBURLShortener/internal/models"
)

type CreateShortURLRequest struct {
	URL      string `json:"url" binding:"required"`
	ShortURL string `json:"short_url"`
}

func validateShortURL(shortURL string) error {
	if len(shortURL) != 6 {
		return fmt.Errorf("short URL must be 6 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(shortURL) {
		return fmt.Errorf("short URL can only contain a-z, A-Z and numbers")
	}

	reserved := []string{"static"} // "analytics", "shorten", "s", "api" - длина != 6, поэтому можно не рассматривать
	for _, word := range reserved {
		if strings.EqualFold(shortURL, word) {
			return fmt.Errorf("unsupported url")
		}
	}

	return nil
}

func (r *CreateShortURLRequest) Validate() error {
	parsed, err := url.Parse(r.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if !(parsed.Scheme == "https" || parsed.Scheme == "http") || parsed.Host == "" {
		return fmt.Errorf("unsupported url")
	}

	if r.ShortURL != "" {
		if err := validateShortURL(r.ShortURL); err != nil {
			return err
		}
	}

	return nil
}

func (r *CreateShortURLRequest) ToModel() (*models.URL, error) {
	r.URL = strings.TrimSpace(r.URL)
	r.ShortURL = strings.TrimSpace(r.ShortURL)
	if err := r.Validate(); err != nil {
		return nil, err
	}

	url := &models.URL{
		ShortURL:    r.ShortURL,
		OriginalURL: r.URL,
	}

	return url, nil

}
