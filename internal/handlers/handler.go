package handlers

import (
	//_ "AvitoTesting/cmd/main/docs"
	// "AvitoTesting/pkg/client/models"
	// "AvitoTesting/pkg/utils"
	"D0SL_organizer/internal/repositories"
	"D0SL_organizer/pkg/client/models"
	"D0SL_organizer/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "github.com/kataras/iris/v12/cache/client"
	// "gorm.io/gorm"
)

type handler struct {
	DB repositories.VideoRepo
}

func New(db repositories.VideoRepo) handler {
	return handler{db}
}

type (
	AddVideoPost struct {
		Link        string `json:"link"`
		Description string `json:"description"`
	}

	GetVideoRequest struct {
		Message string `json:"message"`
	}

	Response struct {
		Message []string `json:"message"`
		Error   string   `json:"error"`
	}
)

func (h handler) AddVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var video models.Video
	err = json.Unmarshal(body, &video)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}
	fmt.Println(video)

	err = h.DB.AddVideo(video)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	}
	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: []string{"Added new video"}})
}

func (h handler) GetSimilarVideos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var embedding GetVideoRequest
	err = json.Unmarshal(body, &embedding)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	vector, err := utils.ParseVector(embedding.Message)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	}

	videos, err := h.DB.GetSimilarVideosByVector(vector)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	}
	utils.RespondWithJSON(w, http.StatusOK, Response{Message: videos})
}
