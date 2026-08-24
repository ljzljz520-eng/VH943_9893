package httpapi

import "net/http"

func (s *Server) inventoryReport(w http.ResponseWriter, r *http.Request) {
	summary, err := s.reporting.Inventory(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
