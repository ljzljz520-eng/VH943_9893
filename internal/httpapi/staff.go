package httpapi

import "net/http"

type staffRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (s *Server) listStaff(w http.ResponseWriter, r *http.Request) {
	staffMembers, err := s.staff.List(r.Context(), r.URL.Query().Get("active") != "false")
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, staffMembers)
}

func (s *Server) enrollStaff(w http.ResponseWriter, r *http.Request) {
	var request staffRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}
	member, err := s.staff.Enroll(r.Context(), request.ID, request.Name, request.Role)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, member)
}
