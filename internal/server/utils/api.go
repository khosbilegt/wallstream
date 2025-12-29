package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

func GetStringFromContext(ctx context.Context, key any) (string, bool) {
	v := ctx.Value(key)
	s, ok := v.(string)
	return s, ok && s != ""
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
