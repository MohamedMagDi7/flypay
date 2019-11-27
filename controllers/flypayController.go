package controllers

import (
	"github.com/callicoder/go-docker/models"
	util "github.com/callicoder/go-docker/utils"
	"github.com/gorilla/schema"
	"net/http"
)

//our API handler for FlyPays
var GetFlyPays = func(w http.ResponseWriter, r *http.Request) {
	// Initialize our Filter Struct
	var filter models.Filter

	//Initialize Our Paramter Decoder
	var paramsDecoder = schema.NewDecoder()

	//Decode The Query Paramters Into filter
	err := paramsDecoder.Decode(&filter, r.URL.Query())

	//check for decoding error
	if err != nil {
		util.Respond(w, util.Message(http.StatusBadRequest, err.Error()))
		return
	}

	//Call the util Prepare filter func to set
	//the statusCode according to Each FlyPay model
	util.PerpareFilter(&filter)

	//make a channel for flypay operations
	ch := make(chan *models.Response)

	//to close the channel after finishing operations
	defer close(ch)

	//Initialize our returned data Model
	var result = []models.FlyPay{}

	//we iterate over all providers supplied in providers map
	//and check if they match the filter and get data accordingly
	for key, value := range util.Providers {
		if filter.Provider == "" || filter.Provider == key {
			//call GetFlyPays for the current provider
			//if matches the filter provider
			go value.GetFlyPays(filter, ch)
			//get result from channel
			res := <-ch
			//check for errors in result
			if res.Err != nil {
				util.Respond(w, util.Message(http.StatusInternalServerError, res.Err.Error()))
				return
			}
			//append provider flypayA data to the final result
			result = append(result, res.List.FlyPay...)

		}
	}

	// Initialize response model
	// set up response status as true and the message as success
	resp := util.Message(http.StatusOK, "success")
	//append result data to the response
	resp["data"] = result
	//call util Respond to encode and write data to http response writer
	util.Respond(w, resp)
}
