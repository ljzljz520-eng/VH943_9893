package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"campgear/internal/catalog"
	"campgear/internal/maintenance"
	"campgear/internal/rental"
	"campgear/internal/reporting"
	"campgear/internal/staff"
)

type Server struct {
	catalog     *catalog.Service
	rental      *rental.Service
	maintenance *maintenance.Service
	reporting   *reporting.Service
	staff       *staff.Service
	serveMux    *http.ServeMux
}

func NewServer(services *rental.Services) *Server {
	server := &Server{catalog: services.Catalog, rental: rental.NewService(services.Repo, services.Catalog), maintenance: maintenance.NewService(services.Repo, services.Catalog), reporting: reporting.NewService(services.Repo), staff: staff.NewService(services.Repo), serveMux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler { return loggingMiddleware(s.serveMux) }

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.Handler())
}

func (s *Server) routes() {
	s.serveMux.HandleFunc("GET /health", s.health)
	s.serveMux.HandleFunc("GET /api/items", s.listItems)
	s.serveMux.HandleFunc("POST /api/items", s.createItem)
	s.serveMux.HandleFunc("GET /api/items/{id}", s.getItem)
	s.serveMux.HandleFunc("PUT /api/items/{id}", s.updateItem)
	s.serveMux.HandleFunc("POST /api/items/{id}/maintenance", s.openMaintenance)
	s.serveMux.HandleFunc("POST /api/maintenance", s.openMaintenance)
	s.serveMux.HandleFunc("POST /api/maintenance/{id}/close", s.closeMaintenance)
	s.serveMux.HandleFunc("GET /api/maintenance", s.listMaintenance)
	s.serveMux.HandleFunc("POST /api/rentals", s.createRental)
	s.serveMux.HandleFunc("GET /api/rentals/{id}", s.getRental)
	s.serveMux.HandleFunc("POST /api/rentals/{id}/activate", s.activateRental)
	s.serveMux.HandleFunc("POST /api/rentals/{id}/return", s.returnRental)
	s.serveMux.HandleFunc("POST /api/rentals/{id}/cancel", s.cancelRental)
	s.serveMux.HandleFunc("GET /api/reports/inventory", s.inventoryReport)
	s.serveMux.HandleFunc("GET /api/staff", s.listStaff)
	s.serveMux.HandleFunc("POST /api/staff", s.enrollStaff)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Campgear-Version", "1")
		next.ServeHTTP(w, r)
	})
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("request body must be valid JSON")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, catalog.ErrItemNotFound) {
		status = http.StatusNotFound
	}
	if strings.Contains(err.Error(), "UNIQUE") {
		status = http.StatusConflict
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
