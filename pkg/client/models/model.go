package models

type Video struct {
	ID          int64        `json:"id"`
	Link        string       `json:"link"`
	Description string       `json:"description"`
	Vector      [768]float32 `json:"vector"`
}
