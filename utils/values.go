package utils

import "github.com/callicoder/go-docker/models"

type FlyPayStatus int

//define statusCode for different models
const (
	Unknown      FlyPayStatus = 0
	Authorised   FlyPayStatus = 1
	Decline      FlyPayStatus = 2
	Refunded     FlyPayStatus = 3
	B_Authorised FlyPayStatus = 100
	B_Decline    FlyPayStatus = 200
	B_Refunded   FlyPayStatus = 300
)

// we list our providers here
//if we need to Add a new provider
//we need to add its filter name as key
//we add the struct type as value
var Providers = map[string]models.FlyPayInterface{
	"flypayA": &models.FlyPay{},
	"flypayB": &models.FlyPayB{},
}

// return the string value of status code enums
func (status FlyPayStatus) String() string {
	switch status {
	case Authorised, B_Authorised:
		return "authorised"
	case Decline, B_Decline:
		return "decline"
	case Refunded, B_Refunded:
		return "refunded"
	default:
		return "unknown"
	}

}
