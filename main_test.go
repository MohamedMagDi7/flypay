package main

import (
	"github.com/callicoder/go-docker/controllers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/payment/transaction", controllers.GetFlyPays).Methods("GET")
	return router
}

func TestGetFlyPays(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/payment/transaction", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

/**************************BenchMark****************************************/
func BenchmarkGetFlyPays(b *testing.B) {
	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest("GET", "/api/payment/transaction", nil)
		response := httptest.NewRecorder()
		Router().ServeHTTP(response, request)
	}
}
