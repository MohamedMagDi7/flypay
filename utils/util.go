package utils

import (
	"encoding/json"
	"github.com/callicoder/go-docker/models"
	"net/http"
)

//Message utility function to set up response model structure
func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

//Respond utility function used to encode data and write to http response writer
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//PrepareFilter used to customize filter data according to each model
func PerpareFilter(filter *models.Filter) {
	if filter.Status != "" {
		switch filter.Status {
		case Authorised.String():
			filter.AStatusCode = int(Authorised)
			filter.BStatusCode = int(B_Authorised)
		case Decline.String():
			filter.AStatusCode = int(Decline)
			filter.BStatusCode = int(B_Decline)
		case Refunded.String():
			filter.AStatusCode = int(Refunded)
			filter.BStatusCode = int(B_Refunded)
		default:
			filter.AStatusCode = int(Unknown)
			filter.BStatusCode = int(Unknown)

		}
	}
}
