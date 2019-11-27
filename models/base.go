package models

const (
	//json data file paths
	ProviderAFilePath            string = "./models/data/flypayA.json"
	ProviderAAlternativeFilePath string = "./data/flypayA.json"
	ProviderBFilePath            string = "./models/data/flypayB.json"
	ProviderBAlternativeFilePath string = "./data/flypayB.json"
	//number of items per goroutine
	ItemsPerRoutine int = 100
)

//Response returned through channels from Models
type Response struct {
	List *FlyPayList
	Err  error
}

//Base FlyPayList returned to the controller
type FlyPayList struct {
	FlyPay []FlyPay `json:"transactions"`
}

//every FlyPay provider must implement GetFLyPays
type FlyPayInterface interface {
	GetFlyPays(filter Filter, ch chan *Response)
	ConvertFlypayToBase() FlyPay
}
