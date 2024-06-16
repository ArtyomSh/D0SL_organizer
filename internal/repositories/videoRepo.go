package repositories

import "D0SL_organizer/pkg/client/models"

type VideoRepo interface {
	AddVideo(video models.Video) error
	GetSimilarVideosByVector(embedding []float32) ([]string, error)
}
