package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	var response []byte
	var err error
	if response, err = json.Marshal(payload); err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ParseVector(vector string) ([]float32, error) {
	result := make([]float32, 0, 768)
	vectorStr := strings.ReplaceAll(vector, "[", "")
	vectorStr = strings.ReplaceAll(vectorStr, "]", "")
	parts := strings.Split(vectorStr, ",")
	if len(parts) != 768 { // dim must be eq vectorDim
		return nil, fmt.Errorf("bad vector dimension")
	}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		v, err := strconv.ParseFloat(part, 32)
		if err != nil {
			continue
		}
		result = append(result, float32(v))
	}
	return result, nil
}
