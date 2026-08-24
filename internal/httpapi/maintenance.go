package httpapi

import (
	"net/http"
)

type maintenanceRequest struct {
	ID         string `json:"id"`
	ItemID     string `json:"item_id"`
	Reason     string `json:"reason"`
	OpenedDate string `json:"opened_date"`
	Technician string `json:"technician"`
}

func (s *Server) openMaintenance(w http.ResponseWriter, r *http.Request) {
	var request maintenanceRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}
	if request.ItemID == "" {
		request.ItemID = r.PathValue("id")
	}
	order, err := s.maintenance.Open(r.Context(), request.ID, request.ItemID, request.Reason, request.OpenedDate, request.Technician)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, order)
}

func (s *Server) closeMaintenance(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Date       string `json:"date"`
		Technician string `json:"technician"`
		Inspection string `json:"inspection"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, err)
		return
	}
	order, err := s.maintenance.Close(r.Context(), r.PathValue("id"), body.Date, body.Technician, body.Inspection)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

func (s *Server) listMaintenance(w http.ResponseWriter, r *http.Request) {
	orders, err := s.maintenance.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}
