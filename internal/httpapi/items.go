package httpapi

import (
	"net/http"

	"campgear/internal/catalog"
)

func (s *Server) listItems(w http.ResponseWriter, r *http.Request) {
	category := catalog.Category(r.URL.Query().Get("category"))
	onlyListed := r.URL.Query().Get("listed") == "true"
	items, err := s.catalog.List(r.Context(), category, onlyListed)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) createItem(w http.ResponseWriter, r *http.Request) {
	var input catalog.ItemInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := s.catalog.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) getItem(w http.ResponseWriter, r *http.Request) {
	item, err := s.catalog.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) updateItem(w http.ResponseWriter, r *http.Request) {
	var input catalog.ItemInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, err)
		return
	}
	item, err := s.catalog.Update(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
