package controllers

import (
	"encoding/json"
	"github.com/callicoder/go-docker/models"
	util "github.com/callicoder/go-docker/utils"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	flypayA string = "flypayA"
	flypayB string = "flypayB"
)

//it returns mock flypayA data
func GetFlyPaysAMock(ch chan *models.Response) {
	var response *models.Response = &models.Response{
		List: &models.FlyPayList{FlyPay: []models.FlyPay{}},
		Err:  nil,
	}

	response.List.FlyPay = append(response.List.FlyPay, models.FlyPay{
		Amount:         1,
		Currency:       "AUD",
		StatusCode:     1,
		OrderReference: "444",
		TransactionID:  "444",
	})
	response.Err = nil
	ch <- response
}

//it returns mock flypayB data
func GetFlyPaysBMock(ch chan *models.Response) {
	var response *models.Response = &models.Response{
		List: &models.FlyPayList{FlyPay: []models.FlyPay{}},
		Err:  nil,
	}
	response.List.FlyPay = append(response.List.FlyPay, models.FlyPay{
		Amount:         1,
		Currency:       "AUD",
		StatusCode:     100,
		OrderReference: "222",
		TransactionID:  "222",
	})
	response.Err = nil
	ch <- response
}

//we are using this mock provider to test our handler logic independently
var GetFlyPaysMock = func(w http.ResponseWriter, r *http.Request) {
	// Initialize our Filter Struct
	var filter models.Filter
	//Initialize Our Paramter Decoder
	var paramsDecoder = schema.NewDecoder()
	//Decode The Query Paramters Into filter
	err := paramsDecoder.Decode(&filter, r.URL.Query())
	//check for decoding error
	if err != nil {
		w.WriteHeader(500)
		return
	}

	//make a channel for flypayA operations
	c1 := make(chan *models.Response)
	//make a channel for flypayB operations
	c2 := make(chan *models.Response)
	var result = []models.FlyPay{}

	if filter.Provider == "" || filter.Provider == flypayA {
		//use goroutine to fetch payflyA mock data
		go GetFlyPaysAMock(c1)
		//get result from channel 1
		res := <-c1
		//check for errors in result
		if res.Err != nil {
			w.WriteHeader(500)
			return
		}
		//append provider flypayA data to the final result
		result = append(result, res.List.FlyPay...)
	}

	//if provider flypayB is selected or none is selected
	if filter.Provider == "" || filter.Provider == flypayB {
		//use goroutine to fetch payflyB mock data
		go GetFlyPaysBMock(c2)
		//get result from channel 1
		res := <-c2
		//check for errors in result
		if res.Err != nil {
			w.WriteHeader(500)
			return
		}
		//append provider flypayB data to the final result
		result = append(result, res.List.FlyPay...)
	}
	// Initialize response model
	// set up response status as true and the message as success
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")

	resp := util.Message(http.StatusOK, "success")
	//append result data to the response
	resp["data"] = result
	//add count of result list
	resp["count"] = int(len(result))

	json.NewEncoder(w).Encode(resp)
}

func TestGetFlyPaysHandler(t *testing.T) {

	r, err := http.NewRequest("GET", "/api/payment/transaction", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFlyPaysMock)
	handler.ServeHTTP(w, r)

	resp := w.Result()
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)
	assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
	assert.Equal(t, 2, int(result["count"].(float64)), "length should equal 2")
}

func TestGetFlyPaysHandlerforProviderA(t *testing.T) {

	r, err := http.NewRequest("GET", "/api/payment/transaction?provider=flypayA", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFlyPaysMock)
	handler.ServeHTTP(w, r)

	resp := w.Result()
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)
	assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
	assert.Equal(t, 1, int(result["count"].(float64)), "length should equal 1")
}

func TestGetFlyPaysHandlerforProviderB(t *testing.T) {

	r, err := http.NewRequest("GET", "/api/payment/transaction?provider=flypayB", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFlyPaysMock)
	handler.ServeHTTP(w, r)

	resp := w.Result()
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)
	assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
	assert.Equal(t, 1, int(result["count"].(float64)), "length should equal 1")
}
