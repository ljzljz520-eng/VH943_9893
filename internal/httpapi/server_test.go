package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"campgear/internal/catalog"
	"campgear/internal/rental"
	"campgear/internal/storage"
)

func TestHTTPWorkflow(t *testing.T) {
	repo, err := storage.Open("file:http-workflow?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	server := NewServer(rental.NewServices(repo))
	createBody, _ := json.Marshal(catalog.ItemInput{ID: "http-tent", SKU: "HT-1", Name: "HTTP Tent", Category: catalog.CategoryTent, DailyRate: 1000, Deposit: 3000, Stock: 2, StorageBin: "H-1"})
	request := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(createBody))
	recorder := httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	rentalBody := `{"id":"http-rental","customer":"Jo","start_date":"2026-08-03","end_date":"2026-08-04","actor":"desk","lines":[{"item_id":"http-tent","quantity":1,"days":1}]}`
	request = httptest.NewRequest(http.MethodPost, "/api/rentals", bytes.NewBufferString(rentalBody))
	recorder = httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("rental status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/api/reports/inventory", nil)
	recorder = httptest.NewRecorder()
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"available_units":1`)) {
		t.Fatalf("report status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if err := context.Background().Err(); err != nil {
		t.Fatal(err)
	}
}
