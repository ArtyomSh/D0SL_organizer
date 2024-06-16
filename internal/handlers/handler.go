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

	// var find models.Segment
	// if err = h.DB.First(&find, "name = ?", segment.Name).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	// 	return
	// }

	// if find.ID != 0 {
	// 	utils.RespondWithJSON(w, http.StatusOK, Response{Message: "segment already exists"})
	// 	return
	// }

	// if result := h.DB.Create(&segment); result.Error != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
	// 	}
	// 	utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	// 	return
	// }

	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: []string{"Added new video"}})
}

// func (h handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()

// 	body, err := io.ReadAll(r.Body)

// 	if err != nil {
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 		return
// 	}

// 	var in GetSegmentName

// 	err = json.Unmarshal(body, &in)
// 	if err != nil {
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 		return
// 	}

// 	var segment models.Segment
// 	if err := h.DB.First(&segment, "name = ?", in.Name).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
// 			return
// 		}
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 		return
// 	}

// 	if err = h.DB.Delete(&segment).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
// 			return
// 		}
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 		return
// 	}

// 	utils.RespondWithJSON(w, http.StatusOK, Response{Message: "Delete segment"})
// }

// func (h handler) AddUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		utils.RespondWithJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
// 		return
// 	}
// 	defer r.Body.Close()

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 		return
// 	}

// 	var in AddUserPost
// 	err = json.Unmarshal(body, &in)
// 	if err != nil {
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 		return
// 	}

// 	var user models.User
// 	if err := h.DB.First(&user, userId).Error; err != nil {
// 		h.DB.Create(&user)
// 		h.DB.First(&user, userId)
// 		utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "create new user"})
// 	}

// 	for _, segmentName := range in.Add {
// 		var segment models.Segment
// 		if err = h.DB.First(&segment, "name = ?", segmentName).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find segment (Name = %s)", segmentName), Error: err.Error()})
// 				return
// 			}
// 			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 			continue
// 		}

// 		err := h.DB.Model(&user).Association("Segments").Append(&segment)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
// 				return
// 			}
// 			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 			return
// 		}
// 		utils.RespondWithJSON(w, http.StatusOK, Response{Message: "add user in new segment"})
// 	}

// 	for _, segmentName := range in.Delete {
// 		var segment models.Segment
// 		if err = h.DB.First(&segment, "name = ?", segmentName).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find segment (Name = %s)", segmentName), Error: err.Error()})
// 				continue
// 			}
// 			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 			continue
// 		}
// 		err := h.DB.Model(&user).Association("Segments").Delete(segment)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
// 				return
// 			}
// 			utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
// 			return
// 		}
// 		utils.RespondWithJSON(w, http.StatusOK, Response{Message: "delete user from segment"})

// 	}

// 	utils.RespondWithJSON(w, http.StatusCreated, Response{Message: "Add/Delete segments from user"})
// }

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

	// var segments []models.Segment
	// var user models.User

	// if err := h.DB.First(&user, userId).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		utils.RespondWithJSON(w, http.StatusNotFound, Response{Message: fmt.Sprintf("cant find user. id=%d", userId), Error: err.Error()})
	// 		return
	// 	}
	// 	utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	// 	return
	// }

	// err = h.DB.Model(&user).Association("Segments").Find(&segments)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		utils.RespondWithJSON(w, http.StatusNotFound, Response{Error: err.Error()})
	// 		return
	// 	}
	// 	utils.RespondWithJSON(w, http.StatusInternalServerError, Response{Error: err.Error()})
	// 	return
	// }

	utils.RespondWithJSON(w, http.StatusOK, Response{Message: videos})
}
