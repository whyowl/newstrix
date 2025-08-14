package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"newstrix/internal/search"
	"strconv"
	"time"
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

func (h *SearchHandler) SearchByFilters(w http.ResponseWriter, r *http.Request) {
	request := search.QueryOption{}

	if query := r.URL.Query().Get("query"); query != "" {
		request.Query = &query
	}
	if source := r.URL.Query().Get("source"); source != "" {
		request.Source = &source
	}
	if from := r.URL.Query().Get("from"); from != "" {
		fromTime, err := time.Parse(time.RFC3339, from)
		if err != nil {
			http.Error(w, "Invalid from date format", http.StatusBadRequest)
			return
		}
		request.From = &fromTime
	}
	if to := r.URL.Query().Get("to"); to != "" {
		toTime, err := time.Parse(time.RFC3339, to)
		if err != nil {
			http.Error(w, "Invalid to date format", http.StatusBadRequest)
			return
		}
		request.To = &toTime
	}
	if keywords := r.URL.Query().Get("keywords"); keywords != "" {
		keywordsList := []string{keywords}
		request.Keywords = &keywordsList
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		request.Limit = limitInt
	} else {
		request.Limit = 0
	}

	results, err := h.service.SearchAdvanced(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, results)
}

func (h *SearchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetByID(r.Context(), &id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if item == nil {
		http.NotFound(w, r)
		return
	}

	respondJSON(w, http.StatusOK, item)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
