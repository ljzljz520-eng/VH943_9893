package httpapi

import (
	"net/http"

	"campgear/internal/rental"
)

type rentalRequest struct {
	ID        string `json:"id"`
	Customer  string `json:"customer"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Lines     []struct {
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
		Days     int    `json:"days"`
	} `json:"lines"`
	Actor string `json:"actor"`
}

func (s *Server) createRental(w http.ResponseWriter, r *http.Request) {
	var request rentalRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}
	cart := rental.NewCart(request.ID, request.Customer, request.StartDate, request.EndDate)
	for _, line := range request.Lines {
		if err := s.rental.AddToCart(r.Context(), cart, line.ItemID, line.Quantity, line.Days); err != nil {
			writeError(w, err)
			return
		}
	}
	record, err := s.rental.CreateRental(r.Context(), *cart, request.Actor)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) getRental(w http.ResponseWriter, r *http.Request) {
	record, err := s.rental.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) activateRental(w http.ResponseWriter, r *http.Request) {
	record, err := s.rental.Activate(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) returnRental(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Date string `json:"date"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, err)
		return
	}
	record, err := s.rental.Return(r.Context(), r.PathValue("id"), body.Date)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) cancelRental(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, err)
		return
	}
	record, err := s.rental.Cancel(r.Context(), r.PathValue("id"), body.Reason)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}
