package handler

import (
	"encoding/json"
	"net/http"
	"newstrix/internal/search"
	"strconv"
)

type SearchHandler struct {
	service *search.SearchEngine
}

func NewSearchHandler(s *search.SearchEngine) *SearchHandler {
	return &SearchHandler{service: s}
}

// GET /search/semantic?query=текст&limit=5
func (h *SearchHandler) SemanticSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	var limit int = 20
	if l := r.URL.Query().Get("limit"); l != "" {
		var err error
		limit, err = strconv.Atoi(l)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	results, err := h.service.SearchBySemanticQuery(r.Context(), query, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, results)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
