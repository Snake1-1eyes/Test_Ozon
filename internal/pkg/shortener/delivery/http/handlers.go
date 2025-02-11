package http

import (
	"net/http"

	"github.com/Snake1-1eyes/Test_Ozon/internal/models"
	shortener "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener"
	"github.com/Snake1-1eyes/Test_Ozon/internal/pkg/utils"
)

type ShortenerHandler struct {
	uc shortener.ShortenerUsecase
}

func NewShortenerHandler(uc shortener.ShortenerUsecase) *ShortenerHandler {
	return &ShortenerHandler{
		uc: uc,
	}
}

func (s *ShortenerHandler) CreateAndSaveShortLink(w http.ResponseWriter, r *http.Request) {
	var req models.Shortener
	err := utils.GetRequestData(r, &req)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error unmarshalling")
		return
	}

	if err := models.ValidateBaseURL(&req); err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "invalid original URL: "+err.Error())
		return
	}

	shortURL, err := s.uc.CreateAndSaveShortLink(req.OriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"shortURL": shortURL,
	}
	err = utils.WriteResponseData(w, response, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ShortenerHandler) GetShortLink(w http.ResponseWriter, r *http.Request) {
	var req models.Shortener
	err := utils.GetRequestData(r, &req)
	if err != nil || req.ShortURL == "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error unmarshalling")
		return
	}

	originalURL, err := s.uc.GetShortLink(req.ShortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]string{
		"originalURL": originalURL,
	}
	err = utils.WriteResponseData(w, response, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
