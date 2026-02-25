package response

import "github.com/Oleska1601/WBURLShortener/internal/models"

type CreateShortURLResponse struct {
	ID       int    `json:"id"`
	ShortURL string `json:"short_url"`
}

func ToCreateShortURLResponse(url *models.URL) *CreateShortURLResponse {
	return &CreateShortURLResponse{
		ID:       url.ID,
		ShortURL: url.ShortURL,
	}
}
